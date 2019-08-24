package main

import (
	"context"
	"fmt"
	"log"

	rsocket "github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

func main() {
	err := rsocket.Receive().
		Resume().
		Fragment(1024).
		Acceptor(func(setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) rsocket.RSocket {
			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					log.Printf("request-response msg: %v", msg)
					return mono.Just(msg)
				}),

				rsocket.FireAndForget(func(msg payload.Payload) {
					log.Printf("fire-and-forget msg: %v", msg)
				}),

				rsocket.RequestStream(func(msg payload.Payload) flux.Flux {
					log.Printf("request-stream msg: %v", msg)

					d := msg.DataUTF8()
					return flux.Create(func(ctx context.Context, s flux.Sink) {
						for i := 0; i < 5; i++ {
							s.Next(payload.NewString(d, fmt.Sprintf("%d", i)))
						}
						s.Complete()
					})
				}),

				rsocket.RequestChannel(func(inputs rx.Publisher) flux.Flux {
					receives := make(chan payload.Payload)

					go func() {
						var count int32
						for range receives {
							count++
							log.Printf("count: %d", count)
						}
					}()

					inputs.(flux.Flux).DoFinally(func(s rx.SignalType) {
						log.Printf("signal type: %v", s)
						close(receives)
					}).Subscribe(context.Background(), rx.OnNext(func(input payload.Payload) {
						log.Printf("subscribe payload: %v", input)
						receives <- input
					}))

					return flux.Create(func(ctx context.Context, s flux.Sink) {
						for i := 0; i < 5; i++ {
							s.Next(payload.NewString("data", fmt.Sprintf("%d_from_server", i)))
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
