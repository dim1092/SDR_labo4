package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

const (
	clientAddress = "127.0.0.1"
	clientPort    = "8091"
)

func startClient(network NetworkConfig) {
	c := make(chan Message)
	udp, err := NewUdp(clientAddress, clientPort)
	if err != nil {
		fmt.Println("Client Error:", err)
		return
	}
	go udp.Listen(c)

	for true {
		fmt.Println("input text to analyse")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Client Error:", err)
			return
		}

		// Sending to all servers
		for _, s := range network.Servers {
			msg := Message{
				ReqCnt,
				clientAddress,
				clientPort,
				input,
				CollectedRes{},
			}
			sendMessage(msg, s.Address, s.ListenPort, udp)
		}

		fmt.Println("enter src id to ask result from - enter 0 to quit")
		_, err = fmt.Scanln(&input)
		if err != nil {
			return
		}
		selectedId, _ := strconv.Atoi(input)
		if selectedId == 0 {
			break
		}
		selectedServ := getServerConfig(selectedId, &network)
		msg := Message{
			ReqResponse,
			clientAddress,
			clientPort,
			"",
			CollectedRes{},
		}
		sendMessage(msg, selectedServ.Address, selectedServ.ListenPort, udp)
		var servMsg Message
		for servMsg.MsgType != Update {
			servMsg = <-c
			if servMsg.MsgType == Other {
				fmt.Println("response:", servMsg.Text)
			}
		}
		fmt.Println("--- Result:")
		for k, v := range servMsg.Collected.Collect {
			fmt.Println(k, " : ", v)
		}

	}
	err = udp.Conn.Close()
	if err != nil {
		return
	}
}

func sendMessage(msg Message, toAddress string, toPort string, conn *UdpConn) {
	resAddr, err := net.ResolveUDPAddr("udp", toAddress+":"+toPort)
	_json, err := json.Marshal(msg)
	err = conn.Send(resAddr, _json)
	if err != nil {
		fmt.Println("Client Error:", err)
		return
	}
}
