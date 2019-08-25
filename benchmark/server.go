package benchmark

import (
	"context"
	"fmt"
	"net/http"

	rsocket "github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
)

// RScoketServer is a simple rsocket server.
func RScoketServer() error {
	return rsocket.Receive().
		Resume().
		Fragment(1024).
		Acceptor(func(setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) rsocket.RSocket {
			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					return mono.Just(msg)
				}),
			)
		}).
		Transport("tcp://127.0.0.1:7878").
		Serve(context.Background())
}

// HTTPServer is a simple HTTP server.
func HTTPServer() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "msg: test")
	})

	return http.ListenAndServe("127.0.0.1:7000", nil)
}
