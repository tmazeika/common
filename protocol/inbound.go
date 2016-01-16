package protocol

import (
	"encoding/gob"
	"fmt"
)

func Inbound(dec *gob.Decoder) (<-chan Signal, <-chan struct{}, <-chan error) {
	sigCh := make(chan Signal)
	typedCh := make(chan struct{})
	errCh := make(chan error)

	go func() {
		defer close(sigCh)
		defer close(typedCh)
		defer close(errCh)

		var inType InboundType
		if err := dec.Decode(&inType); err != nil {
			errCh <- err
			return
		}

		switch inType {
		case SignalInbound:
			var sig Signal
			if err := dec.Decode(&sig); err != nil {
				errCh <- err
				return
			}
			sigCh <- sig
		case TypedInbound:
			typedCh <- struct{}{}
		default:
			errCh <- fmt.Errorf("unknown InboundType 0x%x", inType)
		}
	}()

	return sigCh, typedCh, errCh
}
