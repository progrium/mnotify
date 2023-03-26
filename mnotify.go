package mnotify

import (
	"bufio"
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
	funcs     map[string]func()
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
	c.funcs[fmt.Sprintf("%v", uint64(ptr))] = fn
	fmt.Fprintf(c.stdin, "%v %v\n", uint64(ptr), uint64(size))
}

func New() *Command {
	pid := os.Getpid()
	path, _ := exec.LookPath("mnotify")
	if os.Getenv("MNOTIFY_PATH") != "" {
		path = os.Getenv("MNOTIFY_PATH")
	}
	cmd := exec.Command(path, strconv.Itoa(pid))
	cmd.Stderr = os.Stderr
	rc, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	wc, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	c := &Command{
		Cmd:   cmd,
		funcs: make(map[string]func()),
		stdin: wc,
	}
	go func() {
		scanner := bufio.NewScanner(rc)
		for scanner.Scan() {
			c.mu.Lock()
			fn, ok := c.funcs[scanner.Text()]
			if ok {
				go fn()
			}
			c.mu.Unlock()
		}
	}()
	return c
}
