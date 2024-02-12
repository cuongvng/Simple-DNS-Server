package main


type DNSMessage struct {
	header DNSHeader
	question DNSQuestion
	answer DNSAnswer
}

func (m DNSMessage) serialize() []byte{
	var res = []byte{}
	res = append(res, m.header.serialize()...)
	res = append(res, m.question.serialize()...)
	res = append(res, m.answer.serialize()...)
	return res
}
