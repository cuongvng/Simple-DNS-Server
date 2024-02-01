package main


type DNSMessage struct {
	header DNSHeader
	question DNSQuestion
}

func (m DNSMessage) serialize() []byte{
	var res = []byte{}
	res = append(res, m.header.serialize()...)
	res = append(res, m.question.serialize()...)
	return res
}
