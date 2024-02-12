package main

import (
	"encoding/binary"
)

type DNSAnswer struct {
	NAME string
	TYPE uint16
	CLASS uint16 
	TTL uint32
	RDLength uint16
	RData []byte
}

func (a DNSAnswer) serialize() []byte {
	var res = []byte{}
	res = append(res, encode(a.NAME)...)

	var tmp = make([]byte, 10)
	binary.BigEndian.PutUint16(tmp[0:2], a.TYPE)
	binary.BigEndian.PutUint16(tmp[2:4], a.CLASS)
	binary.BigEndian.PutUint32(tmp[4:8], a.TTL)
	binary.BigEndian.PutUint16(tmp[8:10], a.RDLength)
	tmp = append(tmp, a.RData...)

	res = append(res, tmp...)
	return res
}
