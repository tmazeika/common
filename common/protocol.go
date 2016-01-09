package common

import (
    "fmt"
    "net"
    "bytes"
    "encoding/binary"
)

type Tag byte

// TODO: protocol docs
const (
    Downloader    Tag = 0x00

    Uploader      Tag = 0x01

    // UidAssignment is sent from the puncher to the downloader to indicate the
    // UID of the downloader.
    UidAssignment Tag = 0x02

    // UidRequest is sent from the uploader to the puncher to indicate the
    // UID of the downloader it would like to connect to.
    UidRequest    Tag = 0x03

    // PeerNotFound is sent from the puncher to the uploader to indicate that
    // the requested peer (identified by its uid) could not be found.
    PeerNotFound  Tag = 0x04

    // PeerReady is sent from the puncher to the uploader to indicate that the
    // requested peer (identified by its uid) was found and is ready for further
    // communication. The body contains the external IP address of the peer.
    PeerReady     Tag = 0x05

    // UploaderReady is sent from the puncher to the downloader to indicate that
    // the uploader is ready to send files.
    UploaderReady Tag = 0x06

    // FileName is sent from the uploader to the downloader indicating the name
    // of the file about to be sent.
    FileName      Tag = 0x07

    // FileSize is sent from the uploader to the downloader indicating the size
    // of the file about to be sent.
    FileSize      Tag = 0x08

    // FileHash is sent from the uploader to the downloader indicating the hash
    // of the file about to be sent.
    FileHash      Tag = 0x09

    HashMatch     Tag = 0x0A

    HashMismatch  Tag = 0x0B

    // Error is sent from any peer to another indicating that the sender
    // of this message encountered an error, at the fault of either endpoint.
    // The body contains a string describing the issue.
    Error         Tag = 0x0C

    // Halt is sent from any peer to another indicating that all current
    // connections should be closed and that no future communications will take
    // place. The body contains a message describing the reason.
    Halt          Tag = 0x0D

    // Version is sent from one peer to another indicating the version of
    // itself. The body contains the version number.
    Version       Tag = 0x0E

    // Compatible is sent from one peer to another indicating that it thinks
    // that it is compatible with the peer.
    Compatible    Tag = 0x0F

    // Incompatible is sent from one peer to another indicating that it thinks
    // that it is incompatible with the peer.
    Incompatible  Tag = 0x10
)

type Message struct {
    Tag    Tag
    Length uint16
    Value  []byte
}

func (m *Message) MarshalBinary() (data []byte, err error) {
    const PreValueLen = 3
    const SizeUint16  = 2

    buff := bytes.NewBuffer(make([]byte, PreValueLen + m.Length))

    if err = buff.WriteByte(byte(m.Tag)); err != nil {
        return
    }

    lenBuff := bytes.NewBuffer(make([]byte, SizeUint16))
    binary.BigEndian.PutUint16(lenBuff, m.Length)

    n, err := buff.Write(lenBuff)

    if err != nil {
        return
    }

    if n != cap(lenBuff) {
        return nil, fmt.Errorf("Expected %d byte length, got %d", cap(lenBuff), n)
    }

    n, err = buff.Write(m.Value)

    if err != nil {
        return
    }

    if n != int(m.Length) {
        return nil, fmt.Errorf("Expected %d bytes, got %d", m.Length, n)
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

    closeCh := make(chan int)

    go func() {
        <- closeCh
        close(in.Ch)
        close(out.Ch)
        close(out.Done)
    }()

    go func() {
        defer func() { closeCh <- 0 }()

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
        defer func() { closeCh <- 0 }()

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
