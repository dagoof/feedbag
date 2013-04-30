package feedbag

import (
	"net/http"
	"strings"
	"io"
)

type Pool struct {
	in     chan string
	out    map[chan string]bool
	forget chan chan string
}

func NewPool(c chan string) Pool {
	pool := Pool{c, map[chan string]bool{}, make(chan chan string)}
	go pool.Broadcast()
	return pool
}

func (p Pool) Broadcast() {
	for {
		message := <-p.in
		for target, _ := range p.out {
			target <- message
		}
	}
}

func (p Pool) Add() chan string {
	c := make(chan string)
	p.out[c] = true
	return c
}

func (p Pool) Forget(c chan string) {
	delete(p.out, c)
	close(c)
}

func SharedStream(data chan string, preface ...string) http.HandlerFunc {
	pool := NewPool(data)
	return func(w http.ResponseWriter, r *http.Request) {
		messages := pool.Add()
		defer pool.Forget(messages)

		//w.Header().Set("Content-Type", "text/event-stream")
		for _, part := range preface {
			io.Copy(w, strings.NewReader(part))
		}
		if f, ok := w.(http.Flusher); ok {
			for {
				_, err := io.Copy(w, strings.NewReader(<-messages))
				if err != nil {
					return
				}
				f.Flush()
			}
		}
	}
}
