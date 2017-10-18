package api

import (
	"io"
	"os/exec"
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

func NewGit() *Git {
	return &Git{GitExec: "git"}
}
