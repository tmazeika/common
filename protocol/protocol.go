package protocol

import "crypto/sha256"

type ClientType int

// Client types.
const (
    TargetClient ClientType = iota
    SourceClient
)

type FileInfoMessage struct {
    Name string
    Size int64
    Hash [sha256.Size]byte
}
