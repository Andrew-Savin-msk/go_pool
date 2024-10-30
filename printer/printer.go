package printer

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

	gStop context.CancelFunc

	starter sync.Once
	ctx     context.Context
}

// New initializes new StringHandlers structure
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

// Start starts handling process
func (s *StringHandler) Start() error {
	s.starter.Do(s.run)

	return nil
}

// run creates and starts workers
func (s *StringHandler) run() {
	for s.active < s.limit {
		go s.handler(s.ctx, s.queue, uuid.NewString())
		s.active++
	}
}

// GracefulStop stopping handler gracefully
func (s *StringHandler) GracefulStop() {
	s.gStop()
	close(s.queue)
}

// AddWorkers adding amount of active workers
func (s *StringHandler) AddWorkers(amount int) {
	s.limit += amount

	for s.active < s.limit {
		go s.handler(s.ctx, s.queue, uuid.NewString())
		s.active++
	}
}

// DeleteWorkers limits number of active workers
func (s *StringHandler) DeleteWorkers(amount int) {
	s.limit -= amount
}

// Handle adds task to queue
func (s *StringHandler) Handle(task string) {
	s.queue <- task
}

// handler is a worker that handles messages from queue
func (s *StringHandler) handler(ctx context.Context, work <-chan string, uuid string) {
	for {
		if s.active > s.limit {
			return
		}
		select {
		case str := <-work:
			fmt.Printf("worker_uuid: %s msg: %s\n", uuid, str)
		case <-ctx.Done():
			fmt.Printf("worker %s stopped\n", uuid)
			return
		}
	}
}
