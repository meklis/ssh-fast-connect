package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func LoadConfig(path string, Config *Configuration) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	yamlConfig := string(bytes)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		yamlConfig = strings.ReplaceAll(yamlConfig, fmt.Sprintf("${%v}", pair[0]), pair[1])
	}
	err = yaml.Unmarshal([]byte(yamlConfig), &Config)
	if err != nil {
		return err
	}
	return nil
}
