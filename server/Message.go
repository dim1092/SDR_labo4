package main

type MsgType string

const (
	ReqCnt      MsgType = "count"
	ReqResponse MsgType = "response"
	Update      MsgType = "update"
	Other       MsgType = "other"
)

type Message struct {
	MsgType     MsgType
	FromAddress string
	FromPort    string
	Text        string
	Collected   CollectedRes
}
