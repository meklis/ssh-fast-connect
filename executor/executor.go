package executor

import (
	"fmt"
	"os/exec"
	"ssh-connector/config"
	"strings"
)

type Executor struct {
	config  *config.Configuration
	command string
}

func (e *Executor) SetConfig(configuration *config.Configuration) *Executor {
	e.config = configuration
	return e
}

func (e *Executor) PrepareCommand(server *map[string]string) (*Executor, error) {

	//Get command configuration
	commandName := ""
	for key, value := range *server {
		if key == "command" {
			commandName = value
		}
	}
	if commandName == "" {
		return nil, fmt.Errorf("command is required parameter")
	}

	if _, exist := e.config.Commands[commandName]; !exist {
		return nil, fmt.Errorf("command with name '%s' not found in configuration", commandName)
	}
	command := e.config.Commands[commandName]
	for key, val := range *server {
		command = strings.Replace(command, fmt.Sprintf("%%%s%%", key), val, -1)
	}
	e.command = command

	if strings.Contains(command, "%") {
		return nil, fmt.Errorf("command not complited all variables")
	}
	return e, nil
}

func (e *Executor) Exec() (string, error) {
	//	e.command = strings.Replace(e.command, "\"", "\\\"", -1)
	cmd := exec.Command("bash", "-c", e.command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
