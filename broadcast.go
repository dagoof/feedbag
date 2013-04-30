package feedbag

import (
	"fmt"
	"net/http"
)

type Pool struct {
	in          chan string
	out         map[chan string]bool
	add, remove chan chan string
}

func NewPool(c chan string) *Pool {
	pool := Pool{
		c,
		make(map[chan string]bool),
		make(chan chan string),
		make(chan chan string),
	}
	go pool.accept()
	return &pool
}

// Act as a mutex, by only performing one operation on internal hash at a time
// within this select
func (p *Pool) accept() {
	var cs chan string
	var s string

	for {
		select {
		case cs = <-p.remove:
			delete(p.out, cs)
		case cs = <-p.add:
			p.out[cs] = true
		case s = <-p.in:
			p.broadcast(s)
		}
	}
}

func (p *Pool) broadcast(s string) {
	for target, _ := range p.out {
		target <- s
	}
}

func (p *Pool) Add() chan string {
	c := make(chan string)
	p.add <- c
	return c
}

func (p *Pool) Forget(c chan string) {
	p.remove <- c
	close(c)
}

func SharedStream(data chan string, preface ...string) http.HandlerFunc {
	pool := NewPool(data)
	return func(w http.ResponseWriter, r *http.Request) {
		messages := pool.Add()
		defer pool.Forget(messages)

		w.Header().Set("Content-Type", "text/event-stream")
		for _, part := range preface {
			fmt.Fprintf(w, part)
		}
		if f, ok := w.(http.Flusher); ok {
			for {
				_, err := fmt.Fprintf(w, <-messages)
				if err != nil {
					return
				}
				f.Flush()
			}
		}
	}
}
