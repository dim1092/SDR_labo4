package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type Server struct {
	networkConfig *NetworkConfig
	config        *ServerConfig
	collected     CollectedRes
}

func NewServer(id int, ntwConfig *NetworkConfig) (*Server, error) {
	// Find and set this src's config
	for _, serv := range ntwConfig.Servers {
		if serv.Id == id {
			return &Server{
				ntwConfig,
				&serv,
				CollectedRes{Collect: make(map[string]int)},
			}, nil
		}
	}
	return nil, errors.New("src Id not in config file")

}

func (s *Server) start() {
	// Start listening
	udp, err := NewUdp(s.config.Address, s.config.ListenPort)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c := make(chan Message)
	udp.run = true
	go udp.Listen(c)
	s.handleMessage(c, udp)
}

func (s *Server) spreadResult(udp *UdpConn, c chan Message) {
	for _, srvId := range s.config.Neighbours {
		srv := getServerConfig(srvId, s.networkConfig)
		if srv == nil {
			fmt.Printf("Neighbour src with Id %v not found\n", srvId)

		} else {
			if s.updateNeighbour(srv.Address, srv.ListenPort, udp) != nil {
				fmt.Printf("error sending message to src %v\n", srvId)
			}
		}
	}

	for range s.config.Neighbours {
		// Wait for update and update
		var msg Message
		for msg.MsgType != Update {
			msg = <-c
			if msg.MsgType == ReqCnt || msg.MsgType == ReqResponse {
				// Let the client know the src is still busy
				respMsg := Message{
					Other,
					s.config.Address,
					s.config.ListenPort,
					"Server still working",
					CollectedRes{},
				}
				s.respond(msg, respMsg, udp)
			}
		}

		for k, v := range msg.Collected.Collect {
			s.collected.Collect[k] = v
		}
	}
}

func (s *Server) handleMessage(c chan Message, udp *UdpConn) {

	for true {
		cnt := 1
		msg := <-c
		if msg.MsgType == ReqCnt || msg.MsgType == Update {
			for ; cnt <= s.networkConfig.Depth; cnt++ {
				s.collected.Collect[s.config.Cnt] = cntChar(s.config.Cnt[0], msg.Text)
				s.spreadResult(udp, c)
			}
		} else if msg.MsgType == ReqResponse {
			respMsg := Message{
				Update,
				s.config.Address,
				s.config.ListenPort,
				"",
				s.collected,
			}
			s.respond(msg, respMsg, udp)
		} else if msg.MsgType == Update {
			for k, v := range msg.Collected.Collect {
				s.collected.Collect[k] = v
			}
		}
	}
}

func (s *Server) updateNeighbour(address string, port string, udp *UdpConn) error {
	resAddr, err := net.ResolveUDPAddr("udp", address+":"+port)
	if err != nil {
		return err
	}
	msg, err := json.Marshal(
		Message{Update,
			s.config.Address,
			s.config.ListenPort,
			"",
			s.collected,
		})
	if err != nil {
		return err
	}
	return udp.Send(resAddr, msg)
}

func (s *Server) respond(receivedMsg Message, respondMsg Message, udp *UdpConn) {
	resAddr, err := net.ResolveUDPAddr("udp", receivedMsg.FromAddress+":"+receivedMsg.FromPort)
	if err != nil {
		fmt.Println("Server error: ", err.Error())
		return
	}
	_json, _ := json.Marshal(respondMsg)
	err = udp.Send(resAddr, _json)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
