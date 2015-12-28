package common

const (
    // UidLength is the length of the UID that the puncher server issues.
    UidLength = 16
)

type ProtocolMessage byte

const (
    PuncherPing        ProtocolMessage = 0x00
    PuncherPong        ProtocolMessage = 0x01
    DownloadClientType ProtocolMessage = 0x02
    UploadClientType   ProtocolMessage = 0x03
    ChecksumMismatch   ProtocolMessage = 0x04
    ChecksumMatch      ProtocolMessage = 0x05
)
