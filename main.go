package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	srv := &http.Server{Addr: ":8040"}

	eg.Go(func() error {
		return serverApp(srv)
	})

	eg.Go(func() error {
		<-ctx.Done()
		fmt.Println("Shut down server")
		return srv.Shutdown(ctx)
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-sigs:
				cancel()
			}
		}
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("All clear")
}

func serverApp(server *http.Server) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello, World")
	})
	fmt.Println("App server start, Listen: 8040")
	return server.ListenAndServe()

}
