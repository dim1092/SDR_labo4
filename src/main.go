package main

import "fmt"

func main() {
	var ntwConfig NetworkConfig
	err := loadConfigFromFile(&ntwConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Creating and starting servers
	for _, servConf := range ntwConfig.Servers {
		s, err := NewServer(servConf.Id, &ntwConfig)
		if err != nil {
			fmt.Println(err)
			return
		}
		go s.start()
	}

	startClient(ntwConfig)
}
