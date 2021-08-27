package main

import (
	"flag"
	"fmt"
	"os"
	"ssh-connector/config"
	"ssh-connector/executor"
)

var (
	pathConfig string
	conf       = new(config.Configuration)
)

func init() {
	flag.StringVar(&pathConfig, "c", "config.yml", "Configuration file")
	flag.Parse()
}

func main() {
	//Load configuration
	if err := config.LoadConfig(pathConfig, conf); err != nil {
		panic(err)
	}
	var args = os.Args

	//Run by argument
	if len(args) == 2 {
		serverName := args[1]
		server := conf.GetServerByName(serverName)
		if server == nil {
			fmt.Printf("Server with name %s not found\n", serverName)
			os.Exit(1)
		}
		exec := new(executor.Executor)
		if _, err := exec.SetConfig(conf).PrepareCommand(server); err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		if output, err := exec.Exec(); err != nil {
			fmt.Println(output)
			fmt.Println(err.Error())
			os.Exit(2)
		}
		fmt.Printf("Open connection to '%s'\n", serverName)
		os.Exit(0)
	}
	fmt.Println("Unknown arguments")
}
