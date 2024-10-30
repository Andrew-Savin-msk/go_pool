package workerpool

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

const queueLength = 20

type StringHandler struct {
	limit  int
	active int

	queue chan string

	gStop func()

	starter sync.Once
	ctx     context.Context
}

func New(ctx context.Context, startAmount int) *StringHandler {

	ctx, cancel := context.WithCancel(ctx)

	return &StringHandler{
		limit: startAmount,

		queue: make(chan string, queueLength),

		gStop: cancel,

		starter: sync.Once{},
		ctx:     ctx,
	}
}

func (s *StringHandler) Start() error {
	s.starter.Do(s.run)

	return nil
}

func (s *StringHandler) run() {
	for s.active < s.limit {
		go s.handler(s.ctx, s.queue, uuid.NewString())
		s.active++
	}
}

func (s *StringHandler) GracefulStop() {
	s.gStop()
	close(s.queue)
}

func (s *StringHandler) AddWorkers(amount int) {
	s.limit += amount

	for s.active < s.limit {
		go s.handler(s.ctx, s.queue, uuid.NewString())
		s.active++
	}
}

func (s *StringHandler) DeleteWorkers(amount int) {
	s.limit -= amount
}

func (s *StringHandler) Handle(task string) {
	s.queue <- task
}

func (s *StringHandler) handler(ctx context.Context, work <-chan string, uuid string) {
	for {
		if s.active > s.limit {
			return
		}
		select {
		case str := <-work:
			fmt.Printf("worker_uuid: %s msg: %s\n", uuid, str)
		case <-ctx.Done():
			fmt.Printf("worker %s stopped", uuid)
			return
		}
	}
}
