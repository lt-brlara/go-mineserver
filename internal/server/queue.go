package server

import "github.com/blara/go-mineserver/internal/handle"

type Queue struct {
	Events []handle.Request
	Size   uint8
}

func NewEventQueue() *Queue {
	return &Queue{
		Events: make([]handle.Request, 256),
		Size:   0,
	}
}

func (q *Queue) GetSize() uint8 {
	return q.Size
}

func (q *Queue) Enqueue(req handle.Request) {
	q.Events[q.GetSize()] = req
	q.Size++
}

func (q *Queue) Dequeue() handle.Request {
	req := q.Events[q.GetSize()-1]

	q.Events[q.GetSize()-1] = handle.Request{}
	q.Size--

	return req
}

func (s *Server) AddEvent(req handle.Request) {
	s.eventQueue.Enqueue(req)
}
