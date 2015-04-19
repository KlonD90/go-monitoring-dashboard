package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	ServerAddr,err := net.ResolveUDPAddr("udp",":10001")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
