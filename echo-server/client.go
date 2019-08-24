package main

import (
	"context"
	"fmt"
	"log"

	rsocket "github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
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

	// request-response
	result, err := cli.RequestResponse(payload.NewString("Test RequestAndResponse", "text/plain")).Block(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("request-response:", result)

	// fire-and-forget
	for i := 0; i < 5; i++ {
		data := fmt.Sprintf("Test FireAndForget #%d", i)
		cli.FireAndForget(payload.NewString(data, "text/plain"))
	}

	// request-stream
	done := make(chan struct{})
	_, err = cli.RequestStream(payload.NewString("Test RequestStream", "text/plain")).
		DoFinally(func(s rx.SignalType) {
			log.Printf("s: %v", s)
			close(done)
		}).
		DoOnNext(func(elem payload.Payload) {
			m, _ := elem.MetadataUTF8()
			log.Printf("request-stream: %v", m)
		}).
		BlockLast(ctx)
	<-done

	// request-channel
	send := flux.Create(func(ctx context.Context, s flux.Sink) {
		for i := 0; i < 5; i++ {
			s.Next(payload.NewString("client", fmt.Sprintf("%d", i)))
		}
		s.Complete()
	})

	var seq int
	_, err = cli.RequestChannel(send).
		DoOnNext(func(elem payload.Payload) {
			log.Printf("elem: %v", elem)
			log.Printf("%d_from_server", seq)
			seq++
		}).
		BlockLast(ctx)
}
