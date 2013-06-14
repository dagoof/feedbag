package feedbag

import (
	"fmt"
	"net/http"
)

func drain(c chan string) {
	for _ = range c {
	}
}

func SharedStream(data chan string, preface ...string) http.HandlerFunc {
	pool := NewPool(data)

	return func(w http.ResponseWriter, r *http.Request) {
		messages := pool.Add()

		defer func(c chan string) {
			go drain(messages)
			go pool.Forget(messages)
		}(messages)
		defer func() { recover() }()

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
