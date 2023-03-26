package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/frida/frida-go/frida"
)

var fridaScript = `(function(addr, size) {
  let p = ptr(addr);

  let lastHash = undefined;
  function checkHash() {
    let hash = Checksum.compute("sha1", ArrayBuffer.wrap(p, size));
    if (hash !== lastHash) {
			send(addr);
      lastHash = hash;
    }
  }

  Process.setExceptionHandler((details) => {
    checkHash();

    Memory.protect(p, size, 'rw-');
    setTimeout(() => {
      checkHash();
      Memory.protect(p, size, 'r--');
    }, 1);
    return true;
  })
  Memory.protect(p, size, 'r--');
})(%s, %s);
`

func main() {
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	sess, err := frida.Attach(pid)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		args := strings.Split(scanner.Text(), " ")
		script, err := sess.CreateScript(fmt.Sprintf(fridaScript, args[0], args[1]))
		if err != nil {
			log.Fatal(err)
		}
		script.On("message", func(msg string, data []byte) {
			fmt.Printf("%s\n", args[0])
		})
		script.Load()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
