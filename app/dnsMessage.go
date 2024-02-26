package main


type DNSMessage struct {
	header DNSHeader
	question []DNSQuestion
	answer []DNSAnswer
}

func (m DNSMessage) serialize() []byte{
	var res = []byte{}
	res = append(res, m.header.serialize()...)

	var qdcount = m.header.qdcount
	for i:=0; i<int(qdcount); i++{
		res = append(res, m.question[i].serialize()...)
	}

	var ancount = m.header.ancount
	for i:=0; i<int(ancount); i++{
		res = append(res, m.answer[i].serialize()...)
	}
	
	return res
}
