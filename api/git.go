package api

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

const (
	gitExec = "git"
)

type Git struct {
	GitExec   string
	Cmd       *exec.Cmd
	Path      string
	Env       []string
	Args      string
	Pid       int
	Duration  int
	Errors    []string
	Log       string
	Overwrite bool
	Out       bytes.Buffer
	OutPipe   io.ReadCloser
	Err       bytes.Buffer
	ErrPipe   io.ReadCloser
}

func NewGit(args string) *Git {
	return &Git{GitExec: gitExec, Args: args, Cmd: exec.Command(gitExec, strings.Fields(args)...)}
}

func (g *Git) IsRepo() bool {
	if err := exec.Command(g.GitExec, "rev-parse", "--git-dir").Run(); err != nil {
		return false
	}
	return true
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

func (g *Git) Streamer() error {
	g.OutPipe, _ = g.Cmd.StdoutPipe()
	err := g.Cmd.Start()
	if err != nil {
		return err
	}
	g.Pid = g.Cmd.Process.Pid
	g.Reader(g.Stream(os.Stdout))

	return g.Cmd.Wait()
}

func (g *Git) Run() error {
	g.Cmd.Stdout = &g.Out
	g.Cmd.Stderr = &g.Err
	return g.Cmd.Run()
}
