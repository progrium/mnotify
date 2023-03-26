# mnotify

Tool for getting change notifications on memory addresses. Proof of concept.

To compile mnotify, install `frida-core.h` and `libfrida-core.a` as described [here](https://github.com/frida/frida-go#installation).

A library wraps the mnotify executable, which can be used like this:

```golang
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
```
