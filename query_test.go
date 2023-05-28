package dns_toy

import (
	"net"
	"testing"
)

func TestHeaderToBytes(t *testing.T) {
	header := DNSHeader{
		ID:              0x1314,
		Flags:           0,
		QuestionsCount:  1,
		AnswersCount:    0,
		AuthorityCount:  0,
		AdditionalCount: 0,
	}
	b := HeaderToBytes(header)
	t.Log(b)
}

func TestEncodeDNSName(t *testing.T) {
	name := "google.com"
	b := EncodeDNSName(name)
	want := 15
	if got := len(b); got != want {
		t.Errorf("EncodeDNSName(%q) = %v, want %v", name, got, want)
	}
	t.Log(string(b))
}

func TestQuestionToBytes(t *testing.T) {

}

func TestBuildQuery(t *testing.T) {
	query := BuildQuery("example.com", TYPE_A)
	t.Log(query)
}

func TestDNS(t *testing.T) {
	query := BuildQuery("baidu.com", TYPE_A)
	// 设置服务器地址和端口号
	serverAddr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		t.Error("Error resolving UDP address:", err)
	}

	// 创建UDP连接
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		t.Error("Error connecting to UDP server:", err)
	}
	defer conn.Close()

	// 发送UDP请求
	_, err = conn.Write(query)
	if err != nil {
		t.Error("Error sending UDP message:", err)
	}
	msg := make([]byte, 1024)
	n, err := conn.Read(msg)
	if err != nil {
		t.Error("Error reading UDP response:", err)
	}
	t.Log("UDP message sent successfully! n: ", n)

	packet, err := ParseDnsPacket(msg)
	if err != nil {
		t.Error("Error parsing DNS packet:", err)
	}
	t.Logf("DNS packet: %+v", packet)

	t.Log(IPToString(packet.Answers[0].Data))
}

func TestSendQuery(t *testing.T) {
	packet, err := SendQuery("198.41.0.4:53", "example.com", TYPE_A)
	if err != nil {
		t.Error("Error sending DNS query:", err)
	}
	t.Logf("DNS packet: %+v", packet)
}

func TestResolve(t *testing.T) {
	domain := "example.com"
	ip, err := Resolve(domain, TYPE_A)
	if err != nil {
		t.Error("Error resolving domain:", err)
	}
	t.Logf("Domain %s resolved to '%s'", domain, ip)
}
