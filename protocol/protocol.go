package protocol

import "crypto/sha256"

type ClientType uint8

// Client types.
const (
    ClientDownloader ClientType = iota
    ClientUploader
)

type FileInfoMessage struct {
    Name string
    Size int64
    Hash [sha256.Size]byte
}
