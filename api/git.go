package api

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

type Git struct {
	GitExec   string
	Cmd       *exec.Cmd
	Path      string
	Env       []string
	Args      []string
	Stream    bool
	Success   bool
	Pid       int
	Duration  int
	Errors    []string
	Log       string
	Overwrite bool
	OutPipe   io.ReadCloser
	ErrPipe   io.ReadCloser
}

func NewGit(args string, stream bool) *Git {
	return &Git{GitExec: "git", Args: strings.Fields(args), Stream: stream}
}

func (g *Git) Streamer(l *os.File) (<-chan string, <-chan error) {

	lines := make(chan string)
	errc := make(chan error, 1)

	scanner := bufio.NewScanner(g.OutPipe)

	// Register for interrupts so that we can catch it and immediately
	// return...
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	go func() {
		defer close(lines)
		var err error
		var line string

		errc <- func() error {
			for {
				if !scanner.Scan() {
					err = io.EOF
					break
				}
				line = scanner.Text()
				select {
				case lines <- line:
				}
			}
			return err
		}()

	}()

	return lines, errc
}

func (g *Git) Reader(lines <-chan string, errc <-chan error) {
	go func() {
		for {
			select {
			case line := <-lines:
				fmt.Println(line)
			case <-errc:
				return
			}
		}
	}()

}

func (g *Git) Exec() error {
	g.Cmd = exec.Command(g.GitExec, g.Args...)

	if g.Stream {
		g.OutPipe, _ = g.Cmd.StdoutPipe()
		g.Cmd.Start()
		g.Pid = g.Cmd.Process.Pid
		g.Reader(g.Streamer(os.Stdout))
		g.Cmd.Wait()
	}

	return nil
}
