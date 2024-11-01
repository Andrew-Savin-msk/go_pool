package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Andrew-Savin-msk/go_pool/printer"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sh := printer.New(ctx, 10)

	sh.Start()

	go taskCreator(ctx, sh)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sign := <-stop:
		sh.GracefulStop()
		fmt.Printf("stopping application with signal: %s\n", sign.String())
	}

	cancel()

	fmt.Println("ended work")
}

// taskCreator generate tasks for handler
func taskCreator(ctx context.Context, executor *printer.StringHandler) {
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			executor.Handle(generateRandomString(10))
		}
	}
}

// generateRandomString generate random char sequences
func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
