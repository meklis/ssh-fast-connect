package executor

import (
	"fmt"
	"os/exec"
	"ssh-connector/config"
	"strings"
)

type Executor struct {
	config  *config.Configuration
	profile string
}

func (e *Executor) SetConfig(configuration *config.Configuration) *Executor {
	e.config = configuration
	return e
}

func (e *Executor) PrepareCommand(server *map[string]string) (*Executor, error) {

	//Get command configuration
	commandName := ""
	for key, value := range *server {
		if key == "profile" {
			commandName = value
		}
	}
	if commandName == "" {
		return nil, fmt.Errorf("profile is required parameter")
	}

	if _, exist := e.config.Profiles[commandName]; !exist {
		return nil, fmt.Errorf("profile with name '%s' not found in configuration", commandName)
	}
	command := e.config.Profiles[commandName]
	for key, val := range *server {
		command = strings.Replace(command, fmt.Sprintf("%%%s%%", key), val, -1)
	}
	e.profile = command

	if strings.Contains(command, "%") {
		return nil, fmt.Errorf("profile not complited all variables")
	}
	return e, nil
}

func (e *Executor) Exec() (string, error) {
	//	e.command = strings.Replace(e.command, "\"", "\\\"", -1)
	if strings.Contains(e.profile, "%") {
		return "", fmt.Errorf("profile not complited all variables")
	}
	cmd := exec.Command("bash", "-c", e.profile)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
