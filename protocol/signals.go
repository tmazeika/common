package protocol

import "encoding/gob"

func SignalChannel(dec *gob.Decoder) (<-chan Signal) {
	ch := make(chan Signal)
	errCh := make(chan error)

	go func() {
		var sig Signal
		err := dec.Decode(&sig)

		if err != nil {
			ch <- -1
			return
		}

		ch <- sig
	}()

	return ch, errCh
}
