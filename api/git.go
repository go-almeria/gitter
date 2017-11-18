package api

import (
	"bufio"
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
	Success   bool
	Pid       int
	Duration  int
	Errors    []string
	Log       string
	Overwrite bool
	OutPipe   io.ReadCloser
	ErrPipe   io.ReadCloser
}

func NewGit(args string) *Git {
	return &Git{GitExec: "git", Args: strings.Fields(args)}
}

func (g *Git) Stream(l *os.File) (<-chan string, <-chan error) {

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

func (g *Git) Exec() error {

	g.Cmd = exec.Command(g.GitExec, g.Args...)
	g.OutPipe, _ = g.Cmd.StdoutPipe()

	g.Cmd.Start()
	g.Pid = g.Cmd.Process.Pid
	return nil
}

func (g *Git) Wait() error {
	g.Cmd.Wait()
	return nil
}
