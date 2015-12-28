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
    PuncherEndPing     ProtocolMessage = 0x04
    DownloadClientType ProtocolMessage = 0x05
    UploadClientType   ProtocolMessage = 0x06
    ChecksumMismatch   ProtocolMessage = 0x07
    ChecksumMatch      ProtocolMessage = 0x08
)

func Mtob(msg ProtocolMessage) []byte {
    return []byte{byte(msg)}
}
