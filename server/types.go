package main

type NetworkConfig struct {
	Maxsize int
	Depth   int
	Servers []ServerConfig
}

type ServerConfig struct {
	Id         int
	Address    string
	ListenPort string
	Cnt        string
	Neighbours []int
}

type CollectedRes struct {
	Collect map[string]int
}
