package main

import (
	"fmt"
	"net"
	"encoding/binary"
)

func main() {	
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
		
		var response DNSMessage

		// Parse header
		var headerId = binary.BigEndian.Uint16(buf[:2])
		var opcode = uint16((0b01111000 & buf[2]) >> 3)
		var rd bool = (0b00000001 & buf[3]) & 1 == 1
		var rcode uint16
		if opcode == 0 {
			rcode = 0
		} else {
			rcode = 4
		}
	
		response.header = DNSHeader{
			id: headerId,
			flags: getFlags(true, opcode, false, false, rd, false, 0, rcode),
			qdcount: 1,
			ancount: 1,
			nscount: 0,
			arcount: 0,
		}
		response.question = DNSQuestion{
			NAME: "cuongvng.me",
			TYPE: 1,
			CLASS: 1,
		}
		
		response.answer = DNSAnswer{
			NAME: response.question.NAME,
			TYPE: 1,
			CLASS: 1,
			TTL: 60,
			RDLength: 4,
			RData: []byte{8,8,8,8},
		}
		
		_, err = udpConn.WriteToUDP(response.serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
