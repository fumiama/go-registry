package registry

import (
	"errors"
	"io"
	"net"
	"time"

	tea "github.com/fumiama/gofastTEA"
)

var (
	ErrGetKeyTooLong    = errors.New("get key too long")
	ErrDecAck           = errors.New("decrypt ack error")
	ErrInternalServer   = errors.New("internal server error")
	ErrPermissionDenied = errors.New("permission denied")
	ErrSetKeyTooLong    = errors.New("set key too long")
	ErrSetValTooLong    = errors.New("set val too long")
	ErrUnknownAck       = errors.New("unknown ack error")
	ErrNoSuchKey        = errors.New("no such key")
)

type Regedit struct {
	conn net.Conn
	addr string
	tp   tea.TEA
	ts   *tea.TEA
	seq  byte
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
	p.Put()
	r.seq = 0
	return r.conn.Close()
}

func (r *Regedit) Get(key string) (string, error) {
	if len(key) > 127 {
		return "", ErrGetKeyTooLong
	}
	p := NewCmdPacket(CMDGET, StringToBytes(key), &r.tp)
	defer p.Put()
	r.conn.Write(p.Encrypt(r.seq))
	r.seq++
	err := r.ack(p)
	if err != nil {
		return "", err
	}
	err = p.Decrypt(r.seq)
	if err != nil {
		return "", ErrDecAck
	}
	a := BytesToString(p.data)
	r.seq++
	if a == "erro" {
		return "", ErrInternalServer
	}
	if a == "null" {
		return "", ErrNoSuchKey
	}
	return a, nil
}

func (r *Regedit) Set(key, value string) error {
	if r.ts == nil {
		return ErrPermissionDenied
	}
	if len(key) > 127 {
		return ErrSetKeyTooLong
	}
	if len(value) > 127 {
		return ErrSetValTooLong
	}
	p := NewCmdPacket(CMDSET, StringToBytes(key), r.ts)
	defer p.Put()
	r.conn.Write(p.Encrypt(r.seq))
	r.seq++
	ack := NewCmdPacket(CMDACK, nil, &r.tp)
	defer ack.Put()
	err := r.ack(ack)
	if err != nil {
		return err
	}
	err = ack.Decrypt(r.seq)
	if err != nil {
		return ErrDecAck
	}
	a := BytesToString(ack.data)
	r.seq++
	if a == "erro" {
		return ErrInternalServer
	}
	if a != "data" {
		return ErrUnknownAck
	}
	p.Refresh(CMDDAT, StringToBytes(value), r.ts)
	r.conn.Write(p.Encrypt(r.seq))
	r.seq++
	err = r.ack(ack)
	if err != nil {
		return err
	}
	err = ack.Decrypt(r.seq)
	if err != nil {
		return ErrDecAck
	}
	a = BytesToString(ack.data)
	r.seq++
	if a == "erro" {
		return ErrInternalServer
	}
	if a != "succ" {
		return ErrUnknownAck
	}
	return nil
}

func (r *Regedit) Del(key string) error {
	if r.ts == nil {
		return ErrPermissionDenied
	}
	if len(key) > 127 {
		return ErrGetKeyTooLong
	}
	p := NewCmdPacket(CMDDEL, StringToBytes(key), r.ts)
	defer p.Put()
	r.conn.Write(p.Encrypt(r.seq))
	r.seq++
	ack := NewCmdPacket(CMDACK, nil, &r.tp)
	defer ack.Put()
	err := r.ack(ack)
	if err != nil {
		return err
	}
	err = ack.Decrypt(r.seq)
	if err != nil {
		return ErrDecAck
	}
	a := BytesToString(ack.data)
	r.seq++
	if a == "erro" {
		return ErrInternalServer
	}
	if a == "null" {
		return ErrNoSuchKey
	}
	if a != "succ" {
		return ErrUnknownAck
	}
	return nil
}

func (r *Regedit) ack(c *CmdPacket) error {
	// c.ClearData()
	_, err := io.Copy(c, r.conn)
	return err
}
