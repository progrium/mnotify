package main

import (
	"time"
)

func main() {
	i := 0

	mnotify.Observe(&i, func() {

	})
	defer mnotify.Close()

	for {
		i++
		time.Sleep(2 * time.Second)
	}
}
