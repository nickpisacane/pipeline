package pipeline

import (
	"bytes"
	"os/exec"
	"sync"
)

type Pipeline struct {
	// Commands is a slice of pointers `exec.Cmd` 's
	Commands []*exec.Cmd

	// Gaurd buffers and started
	sync.RWMutex

	// Stdout is a buffer from the last commands stdout
	Stdout bytes.Buffer

	// Stderr is a buffer from the last commands stderr
	Stderr bytes.Buffer

	// Prevent running multiple times
	started bool
}

// New Pipeline
func NewPipeline(cmds ...*exec.Cmd) *Pipeline {
	return &Pipeline{
		Commands: cmds,
		started:  false,
	}
}

// Pipe commands, commands are piped one in to the
// other the order they are appended.
func (p *Pipeline) Append(cmd ...*exec.Cmd) {
	p.Commands = append(p.Commands, cmd...)
}

// Runs the commands sequentially. Each
// commands Stdout is set to the Stdin of
// the following command, if exists. The last
// command is run with the Pipeline's
// buffers.
func (p *Pipeline) Run() error {
	p.Lock()
	defer p.Unlock()

	if p.started {
		return nil
	}

	p.started = true
	length := len(p.Commands) - 1

	for i, cmd := range p.Commands {
		next := i + 1
		if next <= length {
			out, err := cmd.Output()
			if err != nil {
				return err
			}
			var input bytes.Buffer
			input.Write(out)
			p.Commands[next].Stdin = &input
			continue
		}
		cmd.Stdout = &p.Stdout
		cmd.Stderr = &p.Stderr
		cmd.Run()
	}
	return nil
}

// Invokes Run method, Returns the bytes from the Stdout buffer.
func (p *Pipeline) Output() ([]byte, error) {
	if err := p.Run(); err != nil {
		return nil, err
	}
	return p.Stdout.Bytes(), nil
}

// Invokes Run method, returns bytes from Stdout and Stderr.
func (p *Pipeline) CombinedOutput() ([]byte, error) {
	if err := p.Run(); err != nil {
		return nil, err
	}
	if _, err := p.Stdout.Write(p.Stderr.Bytes()); err != nil {
		return nil, err
	}
	return p.Stdout.Bytes(), nil
}
