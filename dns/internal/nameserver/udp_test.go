package nameserver

import (
	"net"
	"testing"

	"golang.org/x/net/dns/dnsmessage"
)

func TestDnsQuestions(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:5300")
	if err != nil {
		t.Fatalf("failed to resolve udp addr: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	var msg = dnsmessage.Message{}
	buf, err := msg.Pack()
	if err != nil {
		t.Fatalf("failed to pack dnsmessage: %v", err)
	}

	_, err = conn.Write(buf)
	if err != nil {
		t.Fatalf("failed to write to udp server: %v", err)
	}

	buf = make([]byte, 512)
	_, _, err = conn.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("failed to read from udp server: %v", err)
	}
	t.Logf("response: %s", string(buf))
	// var resp dnsmessage.Message
	// if err = json.Unmarshal(buf, &resp); err != nil {
	// 	t.Fatalf("json.Unmarshal: %v", err)
	// }

}
