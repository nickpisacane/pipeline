package pipeline

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"sync"
)

type Pipeline struct {
	// Commands is a slice of pointers `exec.Cmd` 's
	Commands []*exec.Cmd

	// Stdin is a ReadWriteBuffer type because it allows
	// the Pipeline struct to be created with out seting
	// the Stdin as an empty buffer, in the case that it
	// is not implemented.
	Stdin ReadWriteBuffer

	// Gaurd buffers and started
	sync.RWMutex

	// Stdout is a buffer from the last commands stdout
	Stdout bytes.Buffer

	// Stderr is a buffer from the last commands stderr
	Stderr bytes.Buffer

	// Prevent running multiple times
	started bool
}

// Implements io.ReadWriter and satifies Truncate method for
// bytes.Buffer.
type ReadWriteBuffer interface {
	io.ReadWriter
	Truncate(int)
}

// New Pipeline
func NewPipeline(cmds ...*exec.Cmd) *Pipeline {
	return &Pipeline{
		Commands: cmds,
		started:  false,
	}
}

// Append commands to the chain.
func (p *Pipeline) Append(cmd ...*exec.Cmd) {
	p.Commands = append(p.Commands, cmd...)
}

// Prepend commands to the chain.
func (p *Pipeline) Prepend(cmds ...*exec.Cmd) {
	p.Commands = append(cmds[:], p.Commands[:]...)
}

// Returns a buffer to write to, the buffer will be piped
// to the Stdin of the first command in the chain when
// run. Returns an error if pipeline has been executed.
func (p *Pipeline) StdinPipe() (io.ReadWriter, error) {
	if p.started {
		return nil, errors.New("Pipeline has already been executed.")
	}
	p.Stdin = bytes.NewBuffer([]byte(""))
	return p.Stdin, nil
}

// Runs commands sequentially, chaining the stdio of the commands.
// If the StdioPipe method was called, the buffer will be piped
// to the Stdin of the first command.
func (p *Pipeline) Run() error {
	p.Lock()
	defer p.Unlock()
	length := len(p.Commands) - 1

	if p.started || length == -1 {
		return nil
	}

	p.started = true

	if p.Stdin != nil {
		p.Commands[0].Stdin = p.Stdin
		defer func() {
			p.Stdin.Truncate(0)
		}()
	}

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
	if _, err := p.Stderr.WriteTo(&p.Stdout); err != nil {
		return nil, err
	}
	return p.Stdout.Bytes(), nil
}
