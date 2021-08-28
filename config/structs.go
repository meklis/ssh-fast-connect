package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
)

type Configuration struct {
	Profiles map[string]string `yaml:"profiles"`
	Groups   []ServerGroup     `yaml:"groups"`
}

type ServerGroup struct {
	Name         string              `yaml:"name"`
	Defaults     map[string]string   `yaml:"defaults"`
	Servers      []map[string]string `yaml:"static_servers"`
	ServerSource *string             `yaml:"servers_source"`
}

func (c *Configuration) LoadServersFromSource() error {
	for key, group := range c.Groups {
		if group.ServerSource == nil {
			continue
		}
		content := ""
		err := error(nil)
		content, err = executeScript(*group.ServerSource)
		if err != nil {
			content, err = readFile(*group.ServerSource)
		}
		if err != nil {
			fmt.Printf("WARNING! Error reading server source in group %s with path %s\n", group.Name, *group.ServerSource)
			continue
		}
		serverData, err := parseData(content)
		if err != nil {
			fmt.Printf("WARNING! Error parse content server source in group %s with path %s\n", group.Name, *group.ServerSource)
			continue
		}
		c.Groups[key].Servers = append(group.Servers, serverData...)

	}

	return nil
}

func (c *Configuration) GetServerByName(name string) *map[string]string {
	serverConfiguration := make(map[string]string)
	var serverFound bool
	for _, group := range c.Groups {
		serverConfiguration = group.Defaults
		for _, server := range group.Servers {
			if _, ok := server["name"]; !ok {
				continue
			}
			if server["name"] == name {
				for k, v := range server {
					serverConfiguration[k] = v
				}
				serverFound = true
				break
			}
		}
	}
	if !serverFound {
		return nil
	}
	return &serverConfiguration
}

func readFile(path string) (string, error) {
	if matched, _ := regexp.MatchString(`^\~\/.*`, path); matched {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		path = strings.Replace(path, "~/", usr.HomeDir+"/", 1)
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func executeScript(path string) (string, error) {
	if matched, _ := regexp.MatchString(`^\~\/.*`, path); matched {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		path = strings.Replace(path, "~/", usr.HomeDir+"/", 1)
	}
	cmd := exec.Command(path)
	cmd.Env = os.Environ()
	output, err := cmd.Output()
	return string(output), err
}

func parseData(config string) ([]map[string]string, error) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		config = strings.ReplaceAll(config, fmt.Sprintf("${%v}", pair[0]), pair[1])
	}
	serverData := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(config), &serverData)
	if err != nil {
		err = yaml.Unmarshal([]byte(config), &serverData)
		if err != nil {
			return nil, err
		}
	}
	servers := make([]map[string]string, 0)
	for _, serv := range serverData {
		servers = append(servers, wrapInterfaceToString(serv))
	}
	return servers, nil
}

func wrapInterfaceToString(server map[string]interface{}) map[string]string {
	data := make(map[string]string)
	for key, val := range server {
		data[key] = fmt.Sprintf("%s", val)
	}
	return data
}
