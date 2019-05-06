package internal

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"path"

	"github.com/DATA-DOG/godog"
)

type commandContext struct {
}

// RegisterCommandContext register execute command steps
func RegisterCommandContext(s *godog.Suite) {
	cc := &commandContext{}

	s.Step(`^I run a command "([^"]*)" with args "([^"]*)"$`, cc.iRunACommand)
}

func (c *commandContext) iRunACommand(command, args string) error {
	cArgs := strings.Split(args, " ")
	var out bytes.Buffer

	// #nosec G204
	cmd := exec.Command(path.Join("..", "..", "bin", command), cArgs...)
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("running command %s. Error: %s\n%s", command, err.Error(), out.String())
	}

	return nil
}
