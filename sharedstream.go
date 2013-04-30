package feedbag

import (
	"fmt"
	"net/http"
)

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
