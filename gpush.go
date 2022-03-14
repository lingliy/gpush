package gpush

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lingliy/gpush/message"
	"google.golang.org/grpc"
)

type client struct {
	conn      *grpc.ClientConn
	rpcClient message.Message_PushClient
	mu        sync.Mutex
	api       map[string]func(jsonStr []byte)
}

func (m *client) registerClient(id string) error {
	c := message.NewMessageClient(m.conn)
	rpcClient, err := c.Push(context.Background(), grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	m.rpcClient = rpcClient
	arr := make([]string, len(m.api))
	i := 0
	for k := range m.api {
		arr[i] = k
		i++
	}
	arrData, err := json.Marshal(arr)
	if err != nil {
		return fmt.Errorf("json marshal failed: %s", err)
	}
	err = rpcClient.Send(&message.Request{Api: "/register", Id: id, Data: arrData})
	if err != nil {
		return fmt.Errorf("send failed: %s", err)
	}
	resp, err := rpcClient.Recv()
	if err != nil {
		return fmt.Errorf("recv failed: %s", err)
	}
	if resp.Api == "/error" && string(resp.Data) == "dup id" {
		return fmt.Errorf("dup id")
	}
	return nil
}

func (m *client) processclientMessage() error {
	req, err := m.rpcClient.Recv()
	if err != nil {
		return err
	}
	f, ok := m.api[req.Api]
	if !ok {
		return nil
	}
	f(req.Data)
	return nil
}

type Server interface {
	RegisterOnline(online func(id string, data []byte) error)
	RegisterOffline(offline func(id string))
	SendCmd(id string, api string, data []byte)

	Push(stream message.Message_PushServer) error
}
type server struct {
	message.UnimplementedMessageServer
	cbCmd   map[string]func(api string, data []byte)
	rw      sync.RWMutex
	online  func(id string, data []byte) error
	offline func(id string)
}

func (s *server) RegisterOnline(online func(id string, data []byte) error) {
	s.online = online
}
func (s *server) RegisterOffline(offline func(id string)) {
	s.offline = offline
}

func (s *server) Push(stream message.Message_PushServer) error {
	response, err := stream.Recv()
	if err != nil {
		return err
	}
	if response.Api != "/register" {
		return fmt.Errorf("no register proto %s", response.Api)
	}

	quit := make(chan struct{})
	cmd := func(api string, data []byte) {
		err := stream.Send(&message.Response{Api: api, Data: data})
		if err != nil {
			<-quit
		}
	}
	id := response.Id
	s.rw.RLock()
	_, ok := s.cbCmd[id]
	s.rw.RUnlock()
	if ok {
		return stream.Send(&message.Response{Api: "/error", Data: []byte("dup id")})
	}
	s.rw.Lock()
	s.cbCmd[id] = cmd
	s.rw.Unlock()
	err = s.online(response.Id, response.Data)
	response.Data = nil
	if err == nil {
		stream.Send(&message.Response{Api: "/ok", Data: []byte("success")})
		select {
		case <-stream.Context().Done():
		case <-quit:
		}
		s.offline(response.Id)
	}
	s.rw.Lock()
	delete(s.cbCmd, response.Id)
	s.rw.Unlock()
	return err
}

func (m *client) Run(id string) error {
	err := m.registerClient(id)
	if err != nil {
		return err
	}
	m.rpcClient.CloseSend()
	for {
		err = m.processclientMessage()
		if err != nil {
			return err
		}
	}
}

func (m *client) RegisterCmd(api string, c func([]byte)) {
	_, ok := m.api[api]
	if ok {
		panic("replay api register")
	}
	m.api[api] = c
}

type Client interface {
	RegisterCmd(api string, c func([]byte))
	Run(id string) error
}

func InitClient(conn *grpc.ClientConn) Client {
	return &client{api: make(map[string]func([]byte)), conn: conn}
}

func InitServer(ss *grpc.Server) Server {
	s := &server{
		cbCmd:   make(map[string]func(api string, data []byte)),
		online:  func(id string, data []byte) error { return nil },
		offline: func(id string) {},
	}
	message.RegisterMessageServer(ss, s)
	return s
}

func (s *server) SendCmd(id string, api string, data []byte) {
	s.rw.RLock()
	f, ok := s.cbCmd[id]
	s.rw.RUnlock()
	if !ok {
		return
	}
	f(api, data)
}
