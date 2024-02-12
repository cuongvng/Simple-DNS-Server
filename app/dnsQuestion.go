package main

import (
	"strings"
	"encoding/hex"
	"encoding/binary"
)

type DNSQuestion struct {
	NAME string
	TYPE uint16
	CLASS uint16
}

func encode(domainName string) []byte {
	var res = []byte{}
	var parts = strings.Split(domainName, ".")
	for _, p := range parts {
		var size = len(p)
		res = append(res, byte(size))
		var encoded = hex.EncodeToString([]byte(p))
		res = append(res, []byte(encoded)...)
	}
	res = append(res, byte(0)) // append a nil byte at last
	return res
}

func (q DNSQuestion) serialize() []byte {
	var res = []byte{}
	res = append(res, encode(q.NAME)...)
	
	var t = make([]byte, 16)
	binary.BigEndian.PutUint16(t, q.TYPE)
	var c = make([]byte, 16)
	binary.BigEndian.PutUint16(c, q.CLASS)

	res = append(res, t...)
	res = append(res, c...)

	return res
}