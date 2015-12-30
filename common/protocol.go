package common
import (
    "os"
    "fmt"
    "net"
)

const (
    // UidLength is the length of the UID that the puncher server issues.
    UidLength = 16
)

// Packet is a description of the data sent from one endpoint to another.
type Packet       byte

// ClientType is a body of ClientType indicating the type of client connecting
// to the puncher.
type ClientType   byte

// Verification is a body of Verification indicating the status of verification
// for the received file (checksum verification).
type Verification byte

const (
    // Ping is sent from the puncher to an uploader or downloader.
    Ping                 Packet       = 0x00

    // Pong is sent from the uploader or downloader to the puncher.
    Pong                 Packet       = 0x01

    // ClientType is sent from the uploader or downloader to the puncher
    // followed by a known length byte indicating the type of client
    // (ClientType).
    ClientType           Packet       = 0x02

    // DownloaderClientType is a body of the ClientType Packet.
    DownloaderClientType ClientType   = 0x00

    // UploaderClientType is a body of the ClientType Packet.
    UploaderClientType   ClientType   = 0x01

    // FileName is sent from the uploader to the downloader indicating the name
    // of the file about to be sent.
    FileName             Packet       = 0x03

    // FileSize is sent from the uploader to the downloader indicating the size
    // of the file about to be sent.
    FileSize             Packet       = 0x04

    // FileHash is sent from the uploader to the downloader indicating the hash
    // of the file about to be sent.
    FileHash             Packet       = 0x05

    // Verification is sent from the downloader to the uploader indicating the
    // status of the hash of the received file (Verification).
    Verification         Packet       = 0x06

    // GoodVerification is the body of the Verification Packet indicating
    // verification has succeeded.
    GoodVerification     Verification = 0x00

    // BadVerification is the body of the Verification Packet indicating
    // verification has failed.
    BadVerification      Verification = 0x01
)

var (
    // bodilessPackets is the set of all Packets that do not have a body.
    bodilessPackets = []Packet{
        Ping,
        Pong,
    }

    // fixedLengthPackets is the map of all Packets that have a fixed length
    // body.
    fixedLengthPackets = map[Packet]uint8{
        ClientType:   1,
        FileSize:     8,  // uint64
        FileHash:     32, // sha256
        Verification: 1,
    }
)

type Message struct {
    packet Packet
    body   []byte
}

func MessageChannel(conn net.Conn) (ch chan Message) {
    ch = make(chan Message)

    go func() {
        for {
            packetBuff := make([]byte, 1)

            if _, err := conn.Read(packetBuff); err != nil {
                fmt.Fprintf(os.Stderr, "Read error for '%s': %s", conn.RemoteAddr(), err)
                return
            }

            packet := Packet(packetBuff[0])

            if ArrayContains(bodilessPackets, packet) {
                ch <- Message{
                    packet: packet,
                }
                continue
            }

            len, known := fixedLengthPackets[packet]

            if ! known {
                lenBuff := make([]byte, 1)

                if _, err := conn.Read(lenBuff); err != nil {
                    fmt.Fprintf(os.Stderr, "Read error for '%s': %s", conn.RemoteAddr(), err)
                    return
                }

                len = uint8(lenBuff[0])
            }

            bodyBuff := make([]byte, len)

            if _, err := conn.Read(bodyBuff); err != nil {
                fmt.Fprintf(os.Stderr, "Read error for '%s': %s", conn.RemoteAddr(), err)
                return
            }

            ch <- Message{
                packet: packet,
                body:   bodyBuff,
            }
        }
    }()

    return
}

func ArrayContains(array []Type, value Type) bool {
    for _, v := range array {
        if v == value {
            return true
        }
    }

    return false
}

func Mtob(msg Packet) []byte {
    return []byte{byte(msg)}
}
