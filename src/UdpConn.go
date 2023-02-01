package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type UdpConn struct {
	Conn    *net.UDPConn
	Address *net.UDPAddr
	run     bool
}

func NewUdp(targetAddress string, targetPort string) (*UdpConn, error) {
	resAddr, err := net.ResolveUDPAddr("udp", targetAddress+":"+targetPort)
	if err != nil {
		return nil, err
	}
	return &UdpConn{nil, resAddr, true}, err
}

func (udp *UdpConn) Start() error {
	if udp.Conn == nil {
		conn, err := net.ListenUDP("udp", udp.Address)
		if err != nil {
			return err
		}
		udp.Conn = conn
	}
	return nil
}

func (udp *UdpConn) Send(targetAddress *net.UDPAddr, msg []byte) error {
	if udp.Conn == nil {
		err := udp.Start()
		if err != nil {
			return err
		}
	}
	_, err := udp.Conn.WriteToUDP(msg, targetAddress)
	return err
}

func (udp *UdpConn) Listen(c chan Message) {
	for udp.run {
		if udp.Conn == nil {
			err := udp.Start()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		buf := make([]byte, 1024)
		_len, _, err := udp.Conn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// Convert to json
		var m Message
		err = json.Unmarshal(buf[:_len], &m)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		c <- m
	}
}
