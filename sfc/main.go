package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"ssh-connector/config"
	"ssh-connector/executor"
)

var (
	pathConfig     string
	conf           = new(config.Configuration)
	printHostnames bool
	profile        string
	list           bool
)

func init() {
	flag.StringVar(&pathConfig, "c", "~/.sfc.conf.yml", "Configuration file")
	flag.StringVar(&profile, "p", "", "Name of profile for usage")
	flag.BoolVar(&list, "l", false, "Print list of hosts")
	flag.BoolVar(&printHostnames, "h", false, "Print hosts")
	flag.Parse()
}

func main() {
	//Load configuration
	if err := config.LoadConfig(pathConfig, conf); err != nil {
		panic(err)
	}
	var args = flag.Args()

	conf.LoadServersFromSource()
	if list {
		for _, group := range conf.Groups {
			fmt.Printf("Group %s, parameters %s\n", group.Name, mapToJson(group.Defaults))

			for _, server := range group.Servers {
				fmt.Println(mapToJson(server))
			}
			fmt.Println("")

		}
		os.Exit(0)
	}

	if printHostnames {
		for _, group := range conf.Groups {
			for _, server := range group.Servers {
				fmt.Println(server["name"])
			}
		}
		os.Exit(0)
	}

	//Work with executor
	if len(args) >= 1 {
		for _, serverName := range args {
			if err := runSshConnect(serverName); err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("Open connection to '%s'\n", serverName)
			}
		}
		os.Exit(1)
	}
	fmt.Println(`SSH fast connect v0.1

Usage: sfc <server name 1> [<server name 2>...] 
`)
}

func runSshConnect(serverName string) error {
	server := conf.GetServerByName(serverName)
	if server == nil {
		return fmt.Errorf("server with name %s not found", serverName)
	}
	exec := new(executor.Executor)
	if profile != "" {
		serv := *server
		serv["profile"] = profile
		server = &serv
	}
	if _, err := exec.SetConfig(conf).PrepareCommand(server); err != nil {
		return err
	}
	if output, err := exec.Exec(); err != nil {
		fmt.Println(output)
		return err
	}
	return nil
}

func mapToJson(data map[string]string) string {
	d, _ := json.Marshal(data)
	return string(d)
}
