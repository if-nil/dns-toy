package dns_toy

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type DNSHeader struct {
	ID              uint16
	Flags           uint16
	QuestionsCount  uint16
	AnswersCount    uint16
	AuthorityCount  uint16
	AdditionalCount uint16
}

type DNSQuestion struct {
	QName  []byte
	QType  uint16
	QClass uint16
}

func HeaderToBytes(header DNSHeader) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, header.ID)
	binary.Write(buf, binary.BigEndian, header.Flags)
	binary.Write(buf, binary.BigEndian, header.QuestionsCount)
	binary.Write(buf, binary.BigEndian, header.AnswersCount)
	binary.Write(buf, binary.BigEndian, header.AuthorityCount)
	binary.Write(buf, binary.BigEndian, header.AdditionalCount)
	return buf.Bytes()
}

func QuestionToBytes(question DNSQuestion) []byte {
	buf := new(bytes.Buffer)
	buf.Write(question.QName)
	binary.Write(buf, binary.BigEndian, question.QType)
	binary.Write(buf, binary.BigEndian, question.QClass)
	return buf.Bytes()
}

func EncodeDNSName(domain string) []byte {
	buf := new(bytes.Buffer)
	labels := strings.Split(domain, ".")
	for _, label := range labels {
		buf.WriteByte(byte(len(label)))
		buf.Write([]byte(label))
	}
	buf.WriteByte(0)
	return buf.Bytes()
}

func BuildQuery(domain string, recordType uint16) []byte {
	name := EncodeDNSName(domain)
	id := rand.Intn(1<<16 - 1)
	header := DNSHeader{
		ID:             uint16(id),
		QuestionsCount: 1,
	}
	question := DNSQuestion{
		QName:  name,
		QType:  recordType,
		QClass: CLASS_IN,
	}
	return append(HeaderToBytes(header), QuestionToBytes(question)...)
}

func SendQuery(addr string, domain string, recordType uint16) (DNSPacket, error) {
	query := BuildQuery(domain, recordType)
	// 设置服务器地址和端口号
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return DNSPacket{}, err
	}

	// 创建UDP连接
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return DNSPacket{}, err
	}
	defer conn.Close()

	// 发送UDP请求
	_, err = conn.Write(query)
	if err != nil {
		return DNSPacket{}, err
	}
	msg := make([]byte, 1024)
	n, err := conn.Read(msg)
	if err != nil {
		return DNSPacket{}, err
	}
	return ParseDnsPacket(msg[:n])
}

func GetAnswer(packet DNSPacket) []byte {
	for _, a := range packet.Answers {
		if a.Type == TYPE_A {
			return a.Data
		}
	}
	return nil
}

func GetNameServerIp(packet DNSPacket) []byte {
	for _, a := range packet.Authority {
		if a.Type == TYPE_A {
			return a.Data
		}
	}
	return nil
}

func GetNameserver(packet DNSPacket) []byte {
	for _, a := range packet.Authority {
		if a.Type == TYPE_NS {
			return a.Data
		}
	}
	return nil
}

func Resolve(domain string, recordType uint16) (string, error) {
	nameServer := "198.41.0.4:53"
	for {
		fmt.Printf("Querying %s from %s\n", domain, nameServer)
		response, err := SendQuery(nameServer, domain, recordType)
		if err != nil {
			return "", err
		}
		if ip := GetAnswer(response); ip != nil {
			return string(ip), nil
		}
		if ip := GetNameServerIp(response); ip != nil {
			nameServer = string(ip) + ":53"
		} else if domain := GetNameserver(response); domain != nil {
			nameServer = string(domain) + ":53"
		} else {
			return "", errors.New("no answer")
		}
	}
}
