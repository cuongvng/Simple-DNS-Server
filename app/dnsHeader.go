package main

import (
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
