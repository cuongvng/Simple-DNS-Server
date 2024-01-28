package main

import (
	"fmt"
	"net"
	"encoding/binary"
)

type DNSHeader struct {
	id uint16
	flags uint16 // qr Opcode Aa Tc Ra Z Rcode
	qdcount uint16
	ancount uint16
	nscount uint16
	arcount uint16
} 
func (h DNSHeader) serialize() []byte{
	var bytes = make([]byte, 12)
	binary.BigEndian.PutUint16(bytes[0:2], h.id)
	binary.BigEndian.PutUint16(bytes[2:4], h.flags)
	binary.BigEndian.PutUint16(bytes[4:6], h.qdcount)
	binary.BigEndian.PutUint16(bytes[6:8], h.ancount)
	binary.BigEndian.PutUint16(bytes[8:10], h.nscount)
	binary.BigEndian.PutUint16(bytes[10:12], h.arcount)

	return bytes
}
func getFlags (QR bool, OPCODE uint16, AA, TC, RD, RA bool, Z, RCODE uint16) uint16{
	var res uint16
	if QR {
		res |= 1 << 15
	}

	res |= OPCODE << 11

	if AA {
		res |= 1 << 10
	}
	if TC {
		res |= 1 << 9
	}
	if RD {
		res |= 1 << 8
	}
	if RA {
		res |= 1 << 7
	}
	
	res |= Z << 4
	res |= RCODE

	return res
}

type DNSMessage struct {
	header DNSHeader
}
func (m DNSMessage) serialize() []byte{
	var res = make([]byte, 12)
	res = m.header.serialize()
	return res
}

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
			qdcount: 0,
			ancount: 0,
			nscount: 0,
			arcount: 0,
		}
	
		_, err = udpConn.WriteToUDP(response.serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
