package main

import (
	"context"
	"log"
	"time"

	rsocket "github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
)

func main() {
	ctx := context.Background()
	cli, err := rsocket.Connect().
		Resume().
		Fragment(1024).
		Transport("tcp://127.0.0.1:7878").
		Start(ctx)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	done := make(chan struct{})
	var su rx.Subscription
	flux := cli.RequestStream(payload.NewString("Test Backpressure", "text/plain"))
	flux.
		DoFinally(func(s rx.SignalType) {
			log.Printf("signal: %s", s)
			close(done)
		}).
		Subscribe(ctx,
			rx.OnSubscribe(func(s rx.Subscription) {
				su = s
				su.Request(1)
			}),
			rx.OnNext(func(msg payload.Payload) {
				log.Printf("msg: %v", msg)
				su.Request(2)
				time.Sleep(2 * time.Second)
			}))
	<-done
}
