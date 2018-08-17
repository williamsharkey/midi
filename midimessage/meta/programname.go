package meta

import (
	"fmt"
	"io"
)

// ProgramName represents a MIDI program name
type ProgramName string

/*
ProgramName

FF 08 length text

This optional event is used to embed the patch/program name that is called up by the immediately subsequent
Bank Select and Program Change messages.

It serves to aid the end user in making an intelligent program choice when using different hardware.

This event may appear anywhere in a track, and there may be multiple occurrences within a track.
*/

// String represents the MIDI program name message as a string (for debugging)
func (p ProgramName) String() string {
	return fmt.Sprintf("%T: %#v", p, p.Text())
}

// Text returns the program name
func (p ProgramName) Text() string {
	return string(p)
}

// Raw returns the raw bytes for the message
func (p ProgramName) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteProgramName),
		Data: []byte(p),
	}).Bytes()
}

func (p ProgramName) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return ProgramName(text), nil
}

func (p ProgramName) meta() {}
