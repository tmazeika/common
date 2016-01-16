package protocol

import (
	"encoding/gob"
	"fmt"
)

func Inbound(dec *gob.Decoder) (exit <-chan struct{}, msg <-chan struct{}, err <-chan error) {
	exitCh := make(chan struct{})
	msgCh := make(chan struct{})
	errCh := make(chan error)

	go func() {
		defer close(exitCh)
		defer close(msgCh)
		defer close(errCh)

		var inType InboundType
		if err := dec.Decode(&inType); err != nil {
			errCh <- err
			return
		}

		switch inType {
		case MessageInbound:
			msgCh <- struct{}{}
		case ExitInbound:
			exitCh <- struct{}{}
		default:
			errCh <- fmt.Errorf("unknown InboundType 0x%x", inType)
		}
	}()

	return exitCh, msgCh, errCh
}
