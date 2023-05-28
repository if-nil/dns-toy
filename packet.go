package dns_toy

import (
	"bytes"
	"net"
)

type DNSPacket struct {
	Header     DNSHeader
	Question   DNSQuestion
	Answers    []DNSRecord
	Authority  []DNSRecord
	Additional []DNSRecord
}

func ParseRecords(length uint16, reader *bytes.Reader) (records []DNSRecord, err error) {
	for i := 0; i < int(length); i++ {
		record, err := ParseRecord(reader)
		if err != nil {
			return records, err
		}
		records = append(records, record)
	}
	return
}

func ParseDnsPacket(data []byte) (packet DNSPacket, err error) {
	reader := bytes.NewReader(data)
	packet.Header = ParseHeader(reader)
	packet.Question, err = ParseQuestion(reader)
	if err != nil {
		return
	}
	packet.Answers, err = ParseRecords(packet.Header.AnswersCount, reader)
	if err != nil {
		return
	}
	packet.Authority, err = ParseRecords(packet.Header.AuthorityCount, reader)
	if err != nil {
		return
	}
	packet.Additional, err = ParseRecords(packet.Header.AdditionalCount, reader)
	return
}

func IPToString(ipBytes []byte) string {
	ip := net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
	return ip.String()
}
