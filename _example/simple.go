package main

import (
	"log"
	"time"

	"github.com/progrium/mnotify"
)

func main() {
	i := 0

	mnotify.Observe(&i, func() {
		log.Println(i)
	})
	defer mnotify.Close()

	for {
		i++
		time.Sleep(2 * time.Second)
	}
}
