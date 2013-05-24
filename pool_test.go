package feedbag

import (
	"testing"
)

func TestPool(t *testing.T) {
	c := make(chan string)
	p := NewPool(c)

	a := p.Add()
	b := p.Add()
	defer p.Forget(a)
	defer p.Forget(b)

	go func(m string) { c <- m }("Hello!")

	if <-a != <-b {
		t.Fatal("The same message should be broadcasted to all clients")
	}
}

