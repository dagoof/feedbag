package feedbag

import (
	"time"
	"net/http"
)

// Write the elapsed time since startup to all connected clients
func ExampleSharedStream() {
	start := time.Now()

	messages := make(chan string)
	fn := func(c chan string) {
		for {
			c <- time.Now().Sub(start).String()
			time.Sleep(time.Second)
		}
	}
	go fn(messages)
	http.Handle("/stream", SharedStream(messages))
}

