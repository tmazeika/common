package common

import "crypto/sha256"

// Client types.
const (
    Downloader byte = 0x00

    Uploader   byte = 0x01
)

type FileInfoMessage struct {
    Name string
    Size uint
    Hash [sha256.Size]byte
}
