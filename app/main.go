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

		var qdcount = binary.BigEndian.Uint16(buf[4:6])
		var ancount = binary.BigEndian.Uint16(buf[6:8])
		var nscount = binary.BigEndian.Uint16(buf[8:10])
		var arcount = binary.BigEndian.Uint16(buf[10:12])

		response.header = DNSHeader{
			id: headerId,
			flags: getFlags(true, opcode, false, false, rd, false, 0, rcode),
			qdcount: qdcount,
			ancount: ancount,
			nscount: nscount,
			arcount: arcount,
		}

		// Parse question: from the 13th byte
		
		var questions []DNSQuestion
		var i = 12
		
		for j:=0; j<int(qdcount); j++ { 
			var i0 = i
			for buf[i] != byte(0) { // get NAME, the nil byte marks the end of it
				i++
			}
			// var qLength = i-i0
			var qName = buf[i0:i]
			var q = DNSQuestion{
				NAME: qName,
				TYPE: 1,
				CLASS: 1,
			}
			questions = append(questions, q)
			i += 1 + 2 + 2 // move to the next question: nil byte (1 bytes), QTYPE (2 bytes), QCLASS (2 bytes)
		}

		response.question = questions
		

		// Parse answer
		var answers []DNSAnswer

		for j:=0; j<int(ancount); j++{
			var i0 = i
			for buf[i] != byte(0) { // get NAME, the nil byte marks the end of it
				i++
			}
			var aName = buf[i0:i]

			i += 2 + 2 + 4 // move over TYPE (2 bytes), CLASS(2 bytes) and TTL (4 bytes)
			var rdLength = binary.BigEndian.Uint16(buf[i:i+2])
			i += 2
			var rData = binary.BigEndian.Uint32(buf[i:i+4])
	
			var a = DNSAnswer{
				NAME: aName,
				TYPE: 1,
				CLASS: 1,
				TTL: 60,
				RDLength: rdLength,
				RData: rData,
			}
			answers = append(answers, a)
		}
		
		response.answer = answers
		
		_, err = udpConn.WriteToUDP(response.serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
