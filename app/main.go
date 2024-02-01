package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}
	
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()
	
	buf := make([]byte, 512)
	
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}
	
		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)
	
		// Create a response with given header
		
		var response DNSMessage
		response.header = DNSHeader{
			id: 1234,
			flags: getFlags(true, 0, false, false, false, false, 0, 0),
			qdcount: 1,
			ancount: 0,
			nscount: 0,
			arcount: 0,
		}
		response.question = DNSQuestion{
			NAME: "codecrafters.io",
			TYPE: 1,
			CLASS: 1,
		}
	
		_, err = udpConn.WriteToUDP(response.serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
