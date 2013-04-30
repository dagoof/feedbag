package main

import (
	"../../feedbag"
	"log"
	"net/http"
	"time"
)

func DoSpam(c chan string, t time.Duration) {
	for {
		c <- "a message"
		time.Sleep(t)
	}
}

func main() {
	spam := make(chan string)
	go DoSpam(spam, time.Second)

	http.Handle("/", feedbag.SharedStream(spam, "preface"))
	log.Fatal(http.ListenAndServe(":8922", nil))
}
