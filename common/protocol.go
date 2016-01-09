package common

import "crypto/sha256"

// Remote types.
const (
    RemoteDownloader = 0
    RemoteUploader   = 1
)

type FileInfoMessage struct {
    Name string
    Size int64
    Hash [sha256.Size]byte
}
