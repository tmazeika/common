package common

const (
    // UidLength is the length of the UID that the puncher server issues.
    UidLength = 16
)

type ProtocolMessage byte

const (
    DownloadClientType ProtocolMessage = 0x00
    UploadClientType   ProtocolMessage = 0x01
    ChecksumMismatch   ProtocolMessage = 0x02
    ChecksumMatch      ProtocolMessage = 0x03
)
