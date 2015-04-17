# pipeline

Provides Unix style piping for the Go's exec Commands. The Pipeline type attempts to stay similar to the os/exec Cmd type.

## Getting Started
```go
pacakge main

import (
	"os/exec"

	"github.com/Nindaff/pipeline"
)

func main() {
	echo := exec.Command("echo", `console.log('print "TEST"');`)
	node := exec.Command("node")
	python := exec.Command("python")
	pipe1 := pipeline.NewPipeline(echo, node, python)
	out, err := pipe1.Output()
	// out => TEST

	// Same command set as pipe1 from string
	pipe2 := pipeline.NewFromString(`echo "console.log('print \'TEST\'');" | node | python`)
}
```
## Usage
### type Pipeline
```go
type Pipeline struct {
	Commands []*exec.Cmd
	Stdin ReadWriteBuffer
	sync.RWMutex
	Stdout bytes.Buffer
	Stderr bytes.Buffer
	started bool
}
```
### ReadWriteBuffer interface
```go
type ReadWriteBuffer interface {
	io.ReadWriter
	Truncate(int)
}
```
### NewPipeline
```go
func NewPipeline(...*exec.Cmd) *Pipeline
```
New Pipeline, optional commands.
### NewFromString
```go
func NewFromString(string) *Pipeline
```
New pipeling from command string.
### ParseString
```go
func ParseString(string) []*exec.Cmd
```
Parsed command string and returns slice of exec.Cmd.

## Pipeline Methods
### Append
```go
func (p *Pipeline) Append(...*exec.Cmd) 
```
Appends commands to the stdio chain.
### Prepend
```go
func (p *Pipeline) Prepend(...*exec.Cmd)
```
Prepend commands to the stdio chain.
### StdinPipe
```go
func (p *Pipeline) StidinPipe() (io.ReadWriter, error)
```
Returns io.ReadWriter, any data in the buffer will be piped to Stdin of first
command.
### Run
```go
func (p *Pipeline) Run() error
```
Runs the commands sequentially.
### Output
```go
func (p *Pipeline) Output() ([]byte, error)
```
Returns bytes from the Stdout.
### CombinedOutput
```go
func (p *Pipeline) CombinedOutput() ([]byte, error)
```
Returns bytes from Stdout and Stderr.

