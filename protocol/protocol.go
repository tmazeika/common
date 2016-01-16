package protocol

type NodeType byte

// NodeType enum.
const (
	TargetNode NodeType = 0x00
	SourceNode NodeType = 0x01
)

type Signal byte

// Signal enum.
const (
	ExitSignal Signal = 0x00
)

type InboundType byte

// InboundType enum.
const (
	SignalInbound InboundType = 0x00
	TypedInbound  InboundType = 0x01
)

type FileInfo struct {
	Name string
	Size uint64
	Hash []byte
}
