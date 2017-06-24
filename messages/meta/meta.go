package meta

import (
	// "fmt"
	"io"
)

/*
from: http://www.somascape.org/midi/tech/mfile.html

Meta events

Meta events are used for special non-MIDI events, and use the 0xFF status that in a MIDI data stream would be used for a System Reset message (a System Reset message would not be useful within a MIDI file).

They have the general form : FF type length data

type specifies the type of Meta event (0 - 127).
length is a variable length quantity (as used to represent delta times) specifying the number of bytes that make up the following data. Some Meta events do not have a data field, whereupon length is 0.

The use of a variable length quantity, rather than a fixed single byte, for length meams that data fields longer than 127 bytes are possible.

The length field should always be read, and should not be assumed, as the definition may change. A MIDI file reader/player should ignore any Meta event types that it does not know about. It should also ignore any additional data if an event's length is longer than expected (it is safe to assume that any extension to the data field will be appended to the current definition). For example if at some time in the future the Sequence Number Meta event is extended with a third data byte, then the first 2 will still have the same interpretation as currently.

Meta event types 0x01 to 0x0F inclusive are reserved for text events. In each case it is best to use the standard 7-bit ASCII character set to ensure reliable interchangeability when transferring files between different computing platforms, however an 8-bit character set may be used. Many text events are best located at or near the beginning of a track (e.g. Copyright, Sequence/Track name, Instrument name), whereas others (Lyric, Marker, Cue point) can occur at various places within a track – their position being an integral aspect of the event.

Although most Meta events are optional, a few are mandatory. Also some events have restrictions regarding their placement.
*/

type Message interface {
	String() string
	Raw() []byte
	meta() // just to tell that it is a meta message
	readFrom(io.Reader) (Message, error)
}

var (
	_ Message = Text("")
	_ Message = Copyright("")
	_ Message = Sequence("")
	_ Message = TrackInstrument("")
	_ Message = Marker("")
	_ Message = Lyric("")
	_ Message = CuePoint("")
	_ Message = SequenceNumber(0)
	_ Message = MIDIChannel(0)
	_ Message = DevicePort("")
	_ Message = MIDIPort(0)
	_ Message = Tempo(0)
	_ Message = TimeSignature{}
	_ Message = KeySignature{}
	_ Message = EndOfTrack
	_ Message = Undefined{}
)

/*
// dunno what it's doing here..

type metaTimeCodeQuarter struct {
	Type   uint8
	Values uint8
}

// TODO check and implement New* function

func (m metaTimeCodeQuarter) String() string {
	return fmt.Sprintf("%#v", m)
}

func (m metaTimeCodeQuarter) meta() {}
*/
