package pipeline

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestPipelineCommands(t *testing.T) {
	assert := assert.New(t)
	cmd1 := exec.Command("echo", `console.log('print "TEST"');`)
	cmd2 := exec.Command("node")
	cmd3 := exec.Command("python")
	pipe := NewPipeline()
	pipe.Append(cmd1, cmd2, cmd3)
	out, err := pipe.Output()
	if err != nil {
		log.Fatal(err)
	}
	strout := strings.Trim(string(out[:]), "\n")
	assert.Equal(strout, "TEST", "Should pipe stdio.")
}

func TestPipelineString(t *testing.T) {
	assert := assert.New(t)
	pipe := NewFromString(`echo "console.log('print \'TEST\'');" | node | python`)
	out, err := pipe.Output()
	if err != nil {
		log.Fatal(err)
	}
	strout := strings.Trim(string(out[:]), "\n")
	assert.Equal(strout, "TEST", "Should parse string commands.")
}
