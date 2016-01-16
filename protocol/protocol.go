package protocol

type NodeType int

// NodeType enum.
const (
	TargetNode = 0
	SourceNode = 1
)

type InboundType int

// InboundType enum.
const (
	MessageInbound = 0
	ExitInbound    = 1
)

type FileInfo struct {
	Name string
	Size uint64
	Hash []byte
}
