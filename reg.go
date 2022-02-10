package registry

import (
	"errors"
	"net"
	"time"

	tea "github.com/fumiama/gofastTEA"
)

type Regedit struct {
	conn net.Conn
	addr string
	tp   tea.TEA
	ts   *tea.TEA
	seq  byte
	buf  [255]byte
}

func NewRegedit(addr, pwd, sps string) *Regedit {
	var tp, ts [16]byte
	if len(pwd) > 15 {
		pwd = pwd[:15]
	}
	if len(sps) > 15 {
		sps = sps[:15]
	}
	copy(tp[:], StringToBytes(pwd))
	copy(ts[:], StringToBytes(sps))
	s := tea.NewTeaCipherLittleEndian(ts[:])
	return &Regedit{addr: addr, tp: tea.NewTeaCipherLittleEndian(tp[:]), ts: &s}
}

func NewRegReader(addr, pwd string) *Regedit {
	var tp [16]byte
	if len(pwd) > 15 {
		pwd = pwd[:15]
	}
	copy(tp[:], StringToBytes(pwd))
	return &Regedit{addr: addr, tp: tea.NewTeaCipherLittleEndian(tp[:])}
}

func (r *Regedit) Connect() (err error) {
	r.conn, err = net.Dial("tcp", r.addr)
	return
}

func (r *Regedit) ConnectIn(wait time.Duration) (err error) {
	r.conn, err = net.DialTimeout("tcp", r.addr, wait)
	return
}

func (r *Regedit) Close() (err error) {
	p := NewCmdPacket(CMDEND, []byte("fill"), &r.tp)
	r.conn.Write(p.Encrypt(r.seq))
	r.seq = 0
	return r.conn.Close()
}

func (r *Regedit) Get(key string) (string, error) {
	if len(key) > 127 {
		return "", errors.New("get key too long")
	}
	p := NewCmdPacket(CMDGET, StringToBytes(key), &r.tp)
	r.conn.Write(p.Encrypt(r.seq))
	r.seq++
	ack, err := r.ack()
	if err != nil {
		return "", err
	}
	ackbytes := ack.Decrypt(r.seq)
	if ackbytes == nil {
		return "", errors.New("decrypt ack error")
	}
	a := BytesToString(ackbytes)
	r.seq++
	if a == "erro" {
		return "", errors.New("server ack error")
	}
	if a == "null" {
		a = ""
	}
	return a, nil
}

func (r *Regedit) Set(key, value string) error {
	if r.ts == nil {
		return errors.New("permission denied")
	}
	if len(key) > 127 {
		return errors.New("set key too long")
	}
	if len(value) > 127 {
		return errors.New("set val too long")
	}
	p := NewCmdPacket(CMDSET, StringToBytes(key), r.ts)
	r.conn.Write(p.Encrypt(r.seq))
	r.seq++
	ack, err := r.ack()
	if err != nil {
		return err
	}
	ackbytes := ack.Decrypt(r.seq)
	if ackbytes == nil {
		return errors.New("decrypt ack error")
	}
	a := BytesToString(ackbytes)
	r.seq++
	if a == "erro" {
		return errors.New("server ack error")
	}
	if a != "data" {
		return errors.New("unknown ack error")
	}
	p = NewCmdPacket(CMDDAT, StringToBytes(value), r.ts)
	r.conn.Write(p.Encrypt(r.seq))
	r.seq++
	ack, err = r.ack()
	if err != nil {
		return err
	}
	ackbytes = ack.Decrypt(r.seq)
	if ackbytes == nil {
		return errors.New("decrypt ack error")
	}
	a = BytesToString(ackbytes)
	r.seq++
	if a == "erro" {
		return errors.New("server ack error")
	}
	if a != "succ" {
		return errors.New("unknown ack error")
	}
	return nil
}

func (r *Regedit) Del(key string) error {
	if r.ts == nil {
		return errors.New("permission denied")
	}
	if len(key) > 127 {
		return errors.New("get key too long")
	}
	p := NewCmdPacket(CMDDEL, StringToBytes(key), r.ts)
	r.conn.Write(p.Encrypt(r.seq))
	r.seq++
	ack, err := r.ack()
	if err != nil {
		return err
	}
	ackbytes := ack.Decrypt(r.seq)
	if ackbytes == nil {
		return errors.New("decrypt ack error")
	}
	a := BytesToString(ackbytes)
	r.seq++
	if a == "erro" {
		return errors.New("server ack error")
	}
	if a == "null" {
		return errors.New("no such key")
	}
	if a != "succ" {
		return errors.New("unknown ack error")
	}
	return nil
}

func (r *Regedit) ack() (*CmdPacket, error) {
	n, err := r.conn.Read(r.buf[:])
	if err != nil {
		return nil, err
	}
	for n < 1+1+16 {
		m, err := r.conn.Read(r.buf[n:])
		if err != nil {
			return nil, err
		}
		n += m
	}
	for n < 1+1+16+int(r.buf[1]) {
		m, err := r.conn.Read(r.buf[n:])
		if err != nil {
			return nil, err
		}
		n += m
	}
	return ParseCmdPacket(r.buf[:], &r.tp), nil
}
