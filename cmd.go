package registry

import (
	"crypto/md5"
	"errors"
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

var (
	ErrMd5Mismatch = errors.New("cmdpacket.decrypt: md5 mismatch")
)

type CmdPacket struct {
	t    tea.TEA
	data []byte
	rawCmdPacket
}

type rawCmdPacket struct {
	cmd uint8
	len uint8
	md5 [16]byte
	raw [255]byte // raw will expand to len
}

//go:nosplit
func NewCmdPacket(cmd uint8, data []byte, t *tea.TEA) *CmdPacket {
	return &CmdPacket{
		t:    *t,
		data: data,
		rawCmdPacket: rawCmdPacket{
			cmd: cmd,
			md5: md5.Sum(data),
		},
	}
}

//go:nosplit
func ParseCmdPacket(data []byte, t *tea.TEA) *CmdPacket {
	if len(data) < 1+1+16 {
		return nil
	}
	if len(data)-1-1-16 < int(data[1]) {
		return nil
	}
	r := (*rawCmdPacket)(*(*unsafe.Pointer)(unsafe.Pointer(&data)))
	c := &CmdPacket{
		t: *t,
		rawCmdPacket: rawCmdPacket{
			cmd: r.cmd,
			len: r.len,
			md5: r.md5,
		},
	}
	copy(c.raw[:], data[1+1+16:])
	return c
}

//go:nosplit
func (c *CmdPacket) Encrypt(seq uint8) (raw []byte) {
	setseq(&c.t, seq)
	c.len = uint8(c.t.EncryptLittleEndianTo(c.data, sumtable, c.raw[:]))
	(*slice)(unsafe.Pointer(&raw)).Data = unsafe.Pointer(&c.rawCmdPacket)
	(*slice)(unsafe.Pointer(&raw)).Len = 1 + 1 + 16 + int(c.len)
	(*slice)(unsafe.Pointer(&raw)).Cap = 1 + 1 + 16 + 255
	return
}

//go:nosplit
func (c *CmdPacket) Decrypt(seq uint8) error {
	setseq(&c.t, seq)
	d := c.t.DecryptLittleEndian(c.raw[:c.len], sumtable)
	if d != nil && c.md5 == md5.Sum(d) {
		c.data = d
		return nil
	}
	return ErrMd5Mismatch
}

//go:nosplit
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
