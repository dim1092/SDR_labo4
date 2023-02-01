package main

type NetworkConfig struct {
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
