package channel

import (
	"fmt"
	"github.com/gomidi/midi/internal/lib"
	"io"
)

const (
	byteProgramChange         = 0xC
	byteChannelPressure       = 0xD
	byteNoteOff               = 0x8
	byteNoteOn                = 0x9
	bytePolyphonicKeyPressure = 0xA
	byteControlChange         = 0xB
	bytePitchWheel            = 0xE
)

// Reader read a channel message
type Reader interface {
	// Read reads a single channel message.
	// It may just be called once per Reader. A second call returns io.EOF
	Read() (Message, error)
}

type ReaderOption func(*reader)

// ReadNoteOffPedantic lets the reader differenciate between "fake" noteoff messages
// (which are in fact noteon messages (typ 9) with velocity of 0) and "real" noteoff messages (typ 8)
// The former are returned as NoteOffPedantic messages and keep the given velocity, the later
// are returned as NoteOff messages without velocity. That means in order to get all noteoff messages,
// there must be checks for NoteOff and NoteOffPedantic (if this option is set).
// If this option is not set, both kinds are returned as NoteOff (default).
func ReadNoteOffPedantic() ReaderOption {
	return func(rd *reader) {
		rd.readNoteOffPedantic = true
	}
}

// NewReader returns a reader that can read a single channel message
// Read may just be called once per Reader. A second call returns io.EOF
func NewReader(input io.Reader, status byte, options ...ReaderOption) Reader {
	rd := &reader{input, status, false, false}

	for _, opt := range options {
		opt(rd)
	}

	return rd
}

type reader struct {
	input               io.Reader
	status              byte
	done                bool
	readNoteOffPedantic bool
}

// Read may just be called once per Reader. A second call returns io.EOF
func (r *reader) Read() (msg Message, err error) {
	if r.done {
		return nil, io.EOF
	}
	var typ, channel, arg1 uint8

	typ, channel = lib.ParseStatus(r.status)

	arg1, err = lib.ReadByte(r.input)
	r.done = true

	if err != nil {
		return
	}

	switch typ {

	// one argument only
	case byteProgramChange, byteChannelPressure:
		msg = r.getMsg1(typ, channel, arg1)

	// two Arguments needed
	default:
		var arg2 byte
		arg2, err = lib.ReadByte(r.input)

		if err != nil {
			return
		}
		msg = r.getMsg2(typ, channel, arg1, arg2)
	}
	return
}

func (r *reader) getMsg1(typ uint8, channel uint8, arg uint8) (msg setter1) {
	switch typ {
	case byteProgramChange:
		msg = ProgramChange{}
	case byteChannelPressure:
		msg = AfterTouch{}
	default:
		panic(fmt.Sprintf("must not happen (typ % X is not an channel message with one argument)", typ))
	}

	msg = msg.set(channel, arg)
	return
}

func (r *reader) getMsg2(typ uint8, channel uint8, arg1 uint8, arg2 uint8) (msg setter2) {

	switch typ {
	case byteNoteOff:
		if r.readNoteOffPedantic {
			msg = NoteOffPedantic{}
		} else {
			msg = NoteOff{}
		}
	case byteNoteOn:
		msg = NoteOn{}
	case bytePolyphonicKeyPressure:
		msg = PolyphonicAfterTouch{}
	case byteControlChange:
		msg = ControlChange{}
	case bytePitchWheel:
		msg = PitchWheel{}
	default:
		panic(fmt.Sprintf("must not happen (typ % X is not an channel message with two arguments)", typ))
	}

	msg = msg.set(channel, arg1, arg2)

	// handle noteOn messages with velocity of 0 as note offs
	if noteOn, is := msg.(NoteOn); is && noteOn.velocity == 0 {
		msg = NoteOff{}
		msg = msg.set(channel, arg1, 0)
	}
	return
}