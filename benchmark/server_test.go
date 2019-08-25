package benchmark

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	rsocket "github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
)

func init() {
	go func() {
		err := RScoketServer()
		if err != nil {
			log.Printf("Err: %s", err)
		} else {
			log.Println("Bye")
		}
	}()

	go func() {
		err := HTTPServer()
		if err != nil {
			log.Printf("Err: %s", err)
		} else {
			log.Println("Bye")
		}
	}()

	time.Sleep(3 * time.Second)
}

func BenchmarkRSocketServer(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		_, err := cli.RequestResponse(payload.NewString("test", "text/plain")).Block(ctx)
		if err != nil {
			b.Logf("err: %s", err)
		}
	}
}

func BenchmarkHTTPServer(b *testing.B) {
	cli := &http.Client{}
	for i := 0; i < b.N; i++ {
		_, err := cli.Get("http://127.0.0.1:7000/")
		if err != nil {
			b.Logf("err: %s", err)
		}
	}
}
