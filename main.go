package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8000"
	}

	destAddr := os.Getenv("RESIZE_ADDR")
	if destAddr == "" {
		destAddr = "worker:8000"
	}

	twilioID := os.Getenv("TWILIO_ID")
	twilioToken := os.Getenv("TWILIO_TOKEN")
	twilioApplication := os.Getenv("TWILIO_APPLICATION")

	if twilioID != "" && twilioToken != "" && twilioApplication != "" {
		log.Println("Configuring Twilio")
		setupTwilio(twilioID, twilioToken, twilioApplication)
	}

	imgurID := os.Getenv("IMGUR_ID")
	if imgurID != "" {
		log.Println("Configuring Imgur")
		setupImgur(imgurID)
	}

	// start http server
	log.Printf("Starting server on: [%s]\n", addr)
	log.Printf("Setting resize server to: [%s]\n", destAddr)
	q, err := serve(addr, destAddr)
	if err != nil {
		log.Printf("unable to start server: %s\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx, done := context.WithTimeout(ctx, time.Second*10)
	defer done()

	quit(ctx, q)
}

func quit(ctx context.Context, fs ...func(context.Context) error) {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	wg := sync.WaitGroup{}

	for _, f := range fs {
		wg.Add(1)

		go func(f func(ctx context.Context) error) {

			err := f(ctx)
			if err != nil {
				log.Printf("did not shutdown cleanly: %s", err)
			}

			wg.Done()
		}(f)
	}

	c := make(chan struct{})

	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		log.Println("clean shutdown")
	case <-ctx.Done():
		log.Println("deadline exceeded, forcing quit")
	}

}
