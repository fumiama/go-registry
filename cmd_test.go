package registry

import (
	"net"
	"testing"

	tea "github.com/fumiama/gofastTEA"
)

func TestCmdPacket(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatal(err)
	}
	tp := tea.NewTeaCipherLittleEndian([]byte("testpwd\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	ts := tea.NewTeaCipherLittleEndian([]byte("testsps\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	var seq byte
	p := NewCmdPacket(CMDGET, []byte("test"), &tp)
	conn.Write(p.Encrypt(seq))
	seq++
	a := string(ack(t, conn, &tp).Decrypt(seq))
	t.Log(a)
	if a != "null" {
		t.Fail()
	}
	seq++
	p = NewCmdPacket(CMDSET, []byte("test"), &ts)
	conn.Write(p.Encrypt(seq))
	seq++
	a = string(ack(t, conn, &tp).Decrypt(seq))
	t.Log(a)
	if a != "data" {
		t.Fail()
	}
	seq++
	p = NewCmdPacket(CMDDAT, []byte("测试"), &ts)
	conn.Write(p.Encrypt(seq))
	seq++
	a = string(ack(t, conn, &tp).Decrypt(seq))
	t.Log(a)
	if a != "succ" {
		t.Fail()
	}
	seq++
	p = NewCmdPacket(CMDGET, []byte("test"), &tp)
	conn.Write(p.Encrypt(seq))
	seq++
	a = string(ack(t, conn, &tp).Decrypt(seq))
	t.Log(a)
	if a != "测试" {
		t.Fail()
	}
	seq++
	p = NewCmdPacket(CMDDEL, []byte("test"), &ts)
	conn.Write(p.Encrypt(seq))
	seq++
	a = string(ack(t, conn, &tp).Decrypt(seq))
	t.Log(a)
	if a != "succ" {
		t.Fail()
	}
	seq++
	p = NewCmdPacket(CMDGET, []byte("test"), &tp)
	conn.Write(p.Encrypt(seq))
	seq++
	a = string(ack(t, conn, &tp).Decrypt(seq))
	t.Log(a)
	if a != "null" {
		t.Fail()
	}
	seq++
}

func ack(t *testing.T, conn net.Conn, tp *tea.TEA) *CmdPacket {
	var buf [1024]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		t.Fatal(err)
	}
	for n < 1+1+16 {
		m, err := conn.Read(buf[n:])
		if err != nil {
			t.Fatal(err)
		}
		n += m
	}
	for n < 1+1+16+int(buf[1]) {
		m, err := conn.Read(buf[n:])
		if err != nil {
			t.Fatal(err)
		}
		n += m
	}
	return ParseCmdPacket(buf[:], tp)
}
