package common

const (
    // UidLength is the length of the UID that the puncher server issues.
    UidLength = 16
)

type ProtocolMessage byte

const (
    PuncherPing        ProtocolMessage = 0x00
    PuncherPong        ProtocolMessage = 0x01
    PuncherReady       ProtocolMessage = 0x02
    PuncherNotReady    ProtocolMessage = 0x03
    DownloadClientType ProtocolMessage = 0x04
    UploadClientType   ProtocolMessage = 0x05
    ChecksumMismatch   ProtocolMessage = 0x06
    ChecksumMatch      ProtocolMessage = 0x07
)

func Mtob(msg ProtocolMessage) []byte {
    return []byte{byte(msg)}
}
