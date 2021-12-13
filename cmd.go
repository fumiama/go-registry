package registry

import (
	"crypto/md5"
	"unsafe"

	tea "github.com/fumiama/gofastTEA"
)

const (
	CMDGET uint8 = iota
	CMDCAT
	CMDMD5
	CMDACK
	CMDEND
	CMDSET
	CMDDEL
	CMDDAT
)

type CmdPacket struct {
	cmd  uint8
	md5  [16]byte
	t    tea.TEA
	data []byte
}

func NewCmdPacket(cmd uint8, data []byte, t *tea.TEA) *CmdPacket {
	return &CmdPacket{
		cmd:  cmd,
		md5:  md5.Sum(data),
		t:    *t,
		data: data,
	}
}

func ParseCmdPacket(data []byte, t *tea.TEA) *CmdPacket {
	if len(data) < 1+1+16 {
		return nil
	}
	if len(data)-1-1-16 < int(data[1]) {
		return nil
	}
	var md5 [16]byte
	copy(md5[:], data[2:18])
	return &CmdPacket{
		cmd:  data[0],
		md5:  md5,
		t:    *t,
		data: data[18 : data[1]+18],
	}
}

func (c *CmdPacket) Encrypt(seq uint8) (raw []byte) {
	setseq(&c.t, seq)
	d := c.t.EncryptLittleEndian(c.data, sumtable)
	raw = append(raw, c.cmd, uint8(len(d)))
	raw = append(raw, c.md5[:]...)
	raw = append(raw, d...)
	return
}

func (c *CmdPacket) Decrypt(seq uint8) []byte {
	setseq(&c.t, seq)
	d := c.t.DecryptLittleEndian(c.data, sumtable)
	if d != nil && c.md5 == md5.Sum(d) {
		return d
	}
	return nil
}

func setseq(t *tea.TEA, seq uint8) {
	*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(t)) + uintptr(15))) = seq
}

// TEA encoding sumtable
var sumtable = [0x10]uint32{
	0x9e3579b9,
	0x3c6ef172,
	0xd2a66d2b,
	0x78dd36e4,
	0x17e5609d,
	0xb54fda56,
	0x5384560f,
	0xf1bb77c8,
	0x8ff24781,
	0x2e4ac13a,
	0xcc653af3,
	0x6a9964ac,
	0x08d12965,
	0xa708081e,
	0x451221d7,
	0xe37793d0,
}
