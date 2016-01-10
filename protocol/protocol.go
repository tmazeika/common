package protocol

import (
    "crypto/sha256"
    "time"
)

type ClientType int
type Signal int

// Client types.
const (
    TargetClient ClientType = iota
    SourceClient
)

// Signals.
const (
    ExitSignal Signal = iota

    // PingSignal is sent from the puncher to a client, expecting a pong in
    // response. The purpose of this is to measure the latency between the
    // endpoints.
    PingSignal

    // PongSignal is sent from a client to the puncher.
    PongSignal

    OkaySignal
)

type SourceReady struct {
    addr   string
    connAt time.Time
}

type FileInfoMessage struct {
    Name string
    Size int64
    Hash [sha256.Size]byte
}
