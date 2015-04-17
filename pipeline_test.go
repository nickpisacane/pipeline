package pipeline

import (
	"log"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func TestPipelineCommands(t *testing.T) {
	assert := assert.New(t)
	cmd1 := exec.Command("echo", `console.log('print "TEST"');`)
	cmd2 := exec.Command("node")
	cmd3 := exec.Command("python")
	pipe := NewPipeline()
	pipe.Append(cmd1, cmd2, cmd3)
	out, err := pipe.Output()
	check(err)
	strout := strings.Trim(string(out[:]), "\n")
	assert.Equal(strout, "TEST", "Should pipe stdio.")
}

func TestPipelineString(t *testing.T) {
	assert := assert.New(t)
	pipe := NewFromString(`echo "console.log('print \'TEST\'');" | node | python`)
	out, err := pipe.Output()
	check(err)
	strout := strings.Trim(string(out[:]), "\n")
	assert.Equal(strout, "TEST", "Should parse string commands.")
}

func TestPipeLineStdinPipe(t *testing.T) {
	assert := assert.New(t)
	cmd1 := exec.Command("node")
	cmd2 := exec.Command("python")
	pipe := NewPipeline(cmd1, cmd2)
	rw, err := pipe.StdinPipe()
	check(err)
	rw.Write([]byte(`console.log('print "TEST"');`))
	out, err := pipe.Output()
	check(err)
	strout := strings.Trim(string(out[:]), "\n")
	assert.Equal(strout, "TEST", "Should pipe data to Stdin from ReadWriter.")
}
