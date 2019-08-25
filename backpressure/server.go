package main

import (
	"context"
	"fmt"
	"log"

	rsocket "github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
)

func main() {
	err := rsocket.Receive().
		Resume().
		Fragment(1024).
		Acceptor(func(setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) rsocket.RSocket {
			return rsocket.NewAbstractSocket(
				rsocket.RequestStream(func(msg payload.Payload) flux.Flux {
					log.Printf("request-stream: %v", msg)

					d := msg.DataUTF8()
					return flux.Create(func(ctx context.Context, s flux.Sink) {
						for i := 0; i < 10; i++ {
							log.Printf("idx: %d", i)
							s.Next(payload.NewString(d, fmt.Sprintf("%d", i)))
						}
						s.Complete()
					})
				}),
			)
		}).
		Transport("tcp://127.0.0.1:7878").
		Serve(context.Background())

	panic(err)
}
