package protocol

import "crypto/sha256"

type ClientType int
type Signal     int

// Client types.
const (
    TargetClient ClientType = iota
    SourceClient
)

// Signals.
const (
    ExitSignal Signal = iota
    OkaySignal
)

type FileInfoMessage struct {
    Name string
    Size int64
    Hash [sha256.Size]byte
}
