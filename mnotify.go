package mnotify

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"sync"
)

func Observe(v any, fn func()) {
	DefaultCommand.Observe(v, fn)
}

func Close() error {
	DefaultCommand.mu.Lock()
	defer DefaultCommand.mu.Unlock()
	DefaultCommand.stdin.Close()
	return DefaultCommand.Cmd.Process.Kill()
}

var DefaultCommand *Command = New()

type Command struct {
	*exec.Cmd
	isRunning bool
	funcs     map[uintptr]func()
	stdin     io.WriteCloser
	mu        sync.Mutex
}

func (c *Command) Observe(v any, fn func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isRunning {
		if err := c.Cmd.Start(); err != nil {
			panic(err)
		}
		c.isRunning = true
	}
	ptr := uintptr(reflect.ValueOf(v).Elem().Addr().UnsafePointer())
	size := reflect.ValueOf(v).Elem().Type().Size()
	c.funcs[ptr] = fn
	fmt.Fprintf(c.stdin, "%v %v\n", uint64(ptr), uint64(size))
}

func New() *Command {
	pid := os.Getpid()
	path, _ := exec.LookPath("mnotify")
	if os.Getenv("MNOTIFY_PATH") != "" {
		path = os.Getenv("MNOTIFY_PATH")
	}
	cmd := exec.Command(path, strconv.Itoa(pid))
	wc, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return &Command{
		Cmd:   cmd,
		funcs: make(map[uintptr]func()),
		stdin: wc,
	}
}
