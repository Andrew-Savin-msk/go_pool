package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	workerpool "github.com/Andrew-Savin-msk/go_pool/wp"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sh := workerpool.New(ctx, 10)

	sh.Start()

	go taskCreator(ctx, sh)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sign := <-stop:
		cancel()

		sh.GracefulStop()

		fmt.Printf("stopping application with signal: %s\n", sign.String())
	}

	fmt.Println("ended work")
}

func taskCreator(ctx context.Context, executor *workerpool.StringHandler) {
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		default:
			executor.Handle(GenerateRandomString(10))
		}
	}
}

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
