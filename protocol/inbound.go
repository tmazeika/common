package protocol

import (
	"encoding/gob"
	"errors"
	"fmt"
)

func Inbound(dec *gob.Decoder) (msg <-chan struct{}, err <-chan error) {
	msgCh := make(chan struct{})
	errCh := make(chan error)

	go func() {
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
			errCh <- errors.New("exiting")
		default:
			errCh <- fmt.Errorf("unknown InboundType 0x%x", inType)
		}
	}()

	return msgCh, errCh
}
