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

		// Parse header: first 12 bytes
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

		// Parse question: from the 13th byte
		// the nil byte marks the end of the label sequence
		var i = 12
		for buf[i] != byte(0) {
			i++
		}
		var qLength = i-12
		var qName = buf[12:12+qLength]

		response.question = DNSQuestion{
			NAME: qName,
			TYPE: 1,
			CLASS: 1,
		}
		
		// Parse answer
		i += 1 + 2 + 2 // move over the question section: nil byte (1 bytes), QTYPE (2 bytes), QCLASS (2 bytes)
		i += qLength + 2 + 2 + 4
		var rdLength = binary.BigEndian.Uint16(buf[i:i+2])
		i += 2
		var rData = binary.BigEndian.Uint32(buf[i:i+4])
		
		response.answer = DNSAnswer{
			NAME: response.question.NAME,
			TYPE: 1,
			CLASS: 1,
			TTL: 60,
			RDLength: rdLength,
			RData: rData,
		}
		
		_, err = udpConn.WriteToUDP(response.serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
