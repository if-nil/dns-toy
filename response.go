package dns_toy

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type DNSRecord struct {
	Name       string
	Type       uint16
	Class      uint16
	TTL        uint32
	DataLength uint16
	Data       []byte
}

func ParseHeader(reader io.Reader) DNSHeader {
	var header DNSHeader
	binary.Read(reader, binary.BigEndian, &header.ID)
	binary.Read(reader, binary.BigEndian, &header.Flags)
	binary.Read(reader, binary.BigEndian, &header.QuestionsCount)
	binary.Read(reader, binary.BigEndian, &header.AnswersCount)
	binary.Read(reader, binary.BigEndian, &header.AuthorityCount)
	binary.Read(reader, binary.BigEndian, &header.AdditionalCount)
	return header
}

func DecodeDNSName(reader *bytes.Reader) []byte {
	var parts [][]byte
	var length uint8
	for {
		if binary.Read(reader, binary.BigEndian, &length); length == 0 {
			break
		}
		if length&0b1100_0000 != 0 {
			parts = append(parts, DecodeCompressedName(length, reader))
			break
		} else {
			label := make([]byte, length)
			binary.Read(reader, binary.BigEndian, &label)
			parts = append(parts, label)
		}
	}
	return bytes.Join(parts, []byte("."))
}

func DecodeCompressedName(length uint8, reader *bytes.Reader) []byte {
	LowPointer := uint16(length&0b0011_1111) << 8
	HighPointer := make([]byte, 1)
	reader.Read(HighPointer)
	pointer := LowPointer | uint16(HighPointer[0])
	currentPos, _ := reader.Seek(0, io.SeekCurrent)
	reader.Seek(int64(pointer), io.SeekStart)
	name := DecodeDNSName(reader)
	reader.Seek(currentPos, io.SeekStart)
	return name
}

func ParseQuestion(reader *bytes.Reader) (DNSQuestion, error) {
	name := DecodeDNSName(reader)
	data := make([]byte, 4)
	n, err := reader.Read(data)
	if err != nil {
		return DNSQuestion{}, err
	}
	if n != 4 {
		return DNSQuestion{}, errors.New("ParseQuestion: read less than 4 bytes")
	}
	return DNSQuestion{
		QName:  name,
		QType:  binary.BigEndian.Uint16(data[:2]),
		QClass: binary.BigEndian.Uint16(data[2:]),
	}, nil
}

func ParseRecord(reader *bytes.Reader) (DNSRecord, error) {
	name := DecodeDNSName(reader)
	data := make([]byte, 10)
	n, err := reader.Read(data)
	if err != nil {
		return DNSRecord{}, err
	}
	if n != 10 {
		return DNSRecord{}, errors.New("ParseRecord: read less than 10 bytes")
	}
	record := DNSRecord{
		Name:       string(name),
		Type:       binary.BigEndian.Uint16(data[:2]),
		Class:      binary.BigEndian.Uint16(data[2:4]),
		TTL:        binary.BigEndian.Uint32(data[4:8]),
		DataLength: binary.BigEndian.Uint16(data[8:]),
	}
	if record.Type == TYPE_NS {
		record.Data = DecodeDNSName(reader)
	} else if record.Type == TYPE_A {
		data = make([]byte, record.DataLength)
		_, err = reader.Read(data)
		if err != nil {
			return DNSRecord{}, err
		}
		record.Data = []byte(IPToString(data))
	} else {
		data = make([]byte, record.DataLength)
		_, err = reader.Read(data)
		if err != nil {
			return DNSRecord{}, err
		}
		record.Data = data
	}
	return record, nil
}
