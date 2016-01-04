package common

import (
    "fmt"
    "net"
    "bytes"
)

type MessageCh struct {
    Ch  chan Message
    Err error
}

type In struct {
    MessageCh
}

type Out struct {
    MessageCh

    Done chan int
}

// Packet is a description of the data sent from one endpoint to another.
type Packet           byte

// TODO: protocol docs
const (
    Downloader    Packet = 0x00

    Uploader      Packet = 0x01

    // UidAssignment is sent from the puncher to the downloader to indicate the
    // UID of the downloader.
    UidAssignment Packet = 0x02

    // UidRequest is sent from the uploader to the puncher to indicate the
    // UID of the downloader it would like to connect to.
    UidRequest    Packet = 0x03

    // PeerNotFound is sent from the puncher to the uploader to indicate that
    // the requested peer (identified by its uid) could not be found.
    PeerNotFound  Packet = 0x04

    // PeerReady is sent from the puncher to the uploader to indicate that the
    // requested peer (identified by its uid) was found and is ready for further
    // communication. The body contains the external IP address of the peer.
    PeerReady     Packet = 0x05

    // UploaderReady is sent from the puncher to the downloader to indicate that
    // the uploader is ready to send files.
    UploaderReady Packet = 0x06

    // FileName is sent from the uploader to the downloader indicating the name
    // of the file about to be sent.
    FileName      Packet = 0x07

    // FileSize is sent from the uploader to the downloader indicating the size
    // of the file about to be sent.
    FileSize      Packet = 0x08

    // FileHash is sent from the uploader to the downloader indicating the hash
    // of the file about to be sent.
    FileHash      Packet = 0x09

    HashMatch     Packet = 0x0A

    HashMismatch  Packet = 0x0B

    // Error is sent from any peer to another indicating that the sender
    // of this message encountered an error, at the fault of either endpoint.
    // The body contains a string describing the issue.
    Error         Packet = 0x0C

    // Halt is sent from any peer to another indicating that all current
    // connections should be closed and that no future communications will take
    // place. The body contains a message describing the reason.
    Halt          Packet = 0x0D

    // Version is sent from one peer to another indicating the version of
    // itself. The body contains the version number.
    Version       Packet = 0x0E

    // Compatible is sent from one peer to another indicating that it thinks
    // that it is compatible with the peer.
    Compatible    Packet = 0x0F

    // Incompatible is sent from one peer to another indicating that it thinks
    // that it is incompatible with the peer.
    Incompatible  Packet = 0x10
)

var (
    // bodilessPackets is the set of all Packets that do not have a body.
    bodilessPackets = []Packet{
        Downloader,
        Uploader,
        PeerNotFound,
        UploaderReady,
        HashMatch,
        HashMismatch,
        Compatible,
        Incompatible,
    }

    // fixedLengthPackets is the map of all Packets that have a fixed length
    // body.
    fixedLengthPackets = map[Packet]uint8{
        FileSize:      8,  // uint64
        FileHash:      32, // sha256
    }
)

// Message is a message from one endpoint to another with a packet and body.
// Some messages may be bodiless, where body will therefore be nil.
type Message struct {
    // Packet is the Packet that describes the body, if present.
    Packet Packet

    // Body is the bytes that the packet describes. May be nil if bodiless.
    Body   []byte
}

func (m Message) MarshalBinary() (data []byte, err error) {
    var buff bytes.Buffer

    buff.WriteByte(byte(m.Packet))

    if ! isBodiless(m.Packet) {
        if _, fixed := fixedLengthPackets[m.Packet]; ! fixed {
            bodyLen := len(m.Body)

            if bodyLen > 0xFF {
                return nil, fmt.Errorf("length of body cannot fit in 1 byte (got %d bytes)", bodyLen)
            }

            buff.WriteByte(byte(len(m.Body)))
        }

        buff.Write(m.Body)
    }

    return buff.Bytes(), nil
}

// MessageChannel returns a 2 channels of Messages for the given Conn. Closes
// both channels upon error or closure.
func MessageChannel(conn net.Conn) (in *In, out *Out) {
    in = &In{
        MessageCh: MessageCh{
            Ch:  make(chan Message),
            Err: nil,
        },
    }
    out = &Out{
        MessageCh: MessageCh{
            Ch:  make(chan Message),
            Err: nil,
        },
        Done: make(chan int),
    }

    close := make(chan int)

    go func() {
        <- close
        close(in.Ch)
        close(out.Ch)
        close(out.Done)
    }()

    go func() {
        defer func() { close <- 0 }()

        for {
            packetBuff := make([]byte, 1)

            if _, err := conn.Read(packetBuff); err != nil {
                in.Err = err
                break
            }

            packet := Packet(packetBuff[0])

            if isBodiless(packet) {
                in.Ch <- Message{ packet, nil }
                continue
            }

            len, known := fixedLengthPackets[packet]

            if ! known {
                lenBuff := make([]byte, 1)

                if _, err := conn.Read(lenBuff); err != nil {
                    in.Err = err
                    break
                }

                len = uint8(lenBuff[0])
            }

            bodyBuff := make([]byte, len)

            if _, err := conn.Read(bodyBuff); err != nil {
                in.Err = err
                break
            }

            in.Ch <- Message{
                Packet: packet,
                Body:   bodyBuff,
            }
        }
    }()

    go func() {
        defer func() { close <- 0 }()

        for {
            data, err := (<- out.Ch).MarshalBinary()

            if err != nil {
                out.Err = err
                break
            }

            if _, err := conn.Write(data); err != nil {
                out.Err = err
                break
            }

            out.Done <- 0
        }
    }()

    return
}

func isBodiless(p Packet) bool {
    for _, v := range bodilessPackets {
        if v == p {
            return true
        }
    }

    return false
}
