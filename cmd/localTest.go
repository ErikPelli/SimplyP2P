package main

import (
	"github.com/ErikPelli/SimplyP2P"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	listenPorts := []string{"2020", "2021", "2022", "2023", "2024", "2025"}
	peers := [][][]string{
		{},
		{{"127.0.0.1", "2020"}},
		{{"127.0.0.1", "2020"}, {"127.0.0.1", "2021"}},
		{{"127.0.0.1", "2020"}, {"127.0.0.1", "2022"}},
		{{"127.0.0.1", "2021"}, {"127.0.0.1", "2023"}},
		{{"127.0.0.1", "2020"}, {"127.0.0.1", "2022"}},
	}

	// Create 6 P2P peers
	for i, port := range listenPorts {
		if _, err := SimplyP2P.NewNode(port, peers[i], &wg); err != nil {
			panic(err)
		}
	}

	wg.Wait()
}