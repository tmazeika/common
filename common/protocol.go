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

type ProtocolMessage byte

const (
    Ping               ProtocolMessage = 0x00
    Pong               ProtocolMessage = 0x01
    PuncherReady       ProtocolMessage = 0x02
    PuncherNotReady    ProtocolMessage = 0x03
    PuncherEndPing     ProtocolMessage = 0x04
    DownloadClientType ProtocolMessage = 0x05
    UploadClientType   ProtocolMessage = 0x06
    ChecksumMismatch   ProtocolMessage = 0x07
    ChecksumMatch      ProtocolMessage = 0x08
)

type Message struct {
    desc ProtocolMessage
    body []byte
}

func MessageChannel(conn net.Conn) (ch chan Message) {
    ch = make(chan Message)

    go func() {
        for {
            descBuff := make([]byte, 1)

            if _, err := conn.Read(descBuff); err != nil {
                fmt.Fprintf(os.Stderr, "Read error for '%s': %s", conn.RemoteAddr(), err)
                return
            }

            desc := ProtocolMessage(descBuff[0])
            lenBuff := make([]byte, 1)

            if _, err := conn.Read(lenBuff); err != nil {
                fmt.Fprintf(os.Stderr, "Read error for '%s': %s", conn.RemoteAddr(), err)
                return
            }

            len := uint8(lenBuff[0])
            bodyBuff := make([]byte, len)

            if _, err := conn.Read(bodyBuff); err != nil {
                fmt.Fprintf(os.Stderr, "Read error for '%s': %s", conn.RemoteAddr(), err)
                return
            }

            ch <- Message{
                desc: desc,
                body: bodyBuff,
            }
        }
    }()

    return
}

func Mtob(msg ProtocolMessage) []byte {
    return []byte{byte(msg)}
}
