package gate

import (
	"encoding/binary"
	"io"
	"reflect"
	"s7/share/net"
	"s7/share/util"

	"github.com/gogo/protobuf/proto"
)

func NewReadWriter() net.MessageReadWriter {
	return &MessageReadWriter{}
}

type MessageReadWriter struct {
}

func (rw *MessageReadWriter) ReadMsgWithLimit(r io.Reader, limit int) (m net.Message, e error) {
	sz := uint32(0)
	e = binary.Read(r, binary.BigEndian, &sz)
	if e != nil {
		return
	}
	if sz >= uint32(limit) {
		e = net.ErrSizeExceedLimit
		return
	}
	raw := make([]byte, sz)
	_, e = io.ReadFull(r, raw)
	if e != nil {
		return
	}

	msg := &Msg{}
	e = proto.Unmarshal(raw, msg)
	util.PanicOnErr(e)

	typ := proto.MessageType(msg.Name)
	v := reflect.New(typ.Elem())

	m = v.Interface().(net.Message)
	e = proto.Unmarshal(msg.Raw, m.(proto.Message))
	util.PanicOnErr(e)

	return
}
func (rw *MessageReadWriter) WriteMsg(w io.Writer, m net.Message) error {
	pb := m.(proto.Message)
	r, e := proto.Marshal(pb)
	util.PanicOnErr(e)

	msg := &Msg{
		Name: proto.MessageName(pb),
		Raw:  r,
	}

	raw, e := proto.Marshal(msg)
	util.PanicOnErr(e)

	sz := uint32(len(raw))

	e = binary.Write(w, binary.BigEndian, sz)
	if e != nil {
		return e
	}

	_, e = w.Write(raw)

	return e
}
