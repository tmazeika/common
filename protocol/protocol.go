package protocol

type NodeType int

// NodeType enum.
const (
	TargetNode = 0
	SourceNode = 1
)

type Signal int

// Signal enum.
const (
	TargetNotFoundSignal = 0
	OkaySignal           = 1
)

type FileInfo struct {
	Name string
	Size int64
	Hash []byte
}
