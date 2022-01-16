package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v2"
)

type Command struct {
	Name    string
	Command []string
}

type Config struct {
	Init     []string
	Commands []Command
}

type Unmarshaller func([]byte, interface{}) error

var (
	locations = []string{
		"/starter",
		"./starter",
		"/etc/starter",
	}
	keepRunning *int32
)

func logF(f string, v ...interface{}) {
	msg := fmt.Sprintf(f, v...)
	fmt.Println("[starter]", msg)
}

func loadConfig(tryFiles []string) (*Config, error) {
	config := &Config{}
	triedPaths := []string{}
	var um Unmarshaller

	for _, fPath := range tryFiles {
		switch fPath[len(fPath)-4:] {
		case "json":
			um = json.Unmarshal
		case "yaml":
			um = json.Unmarshal
		default:
			return nil, fmt.Errorf("unsupported filetype %s", fPath[len(fPath)-4:])
		}
		// try json
		triedPaths = append(triedPaths, fPath)
		if data, err := ioutil.ReadFile(fPath); err == nil {
			err := um(data, config)
			if err != nil {
				logF("unable to decode %s: %s\n", fPath, err)
				continue
			}
			return config, nil
		}
	}
	return nil, fmt.Errorf("unable to find config in %s", strings.Join(triedPaths, ":"))
}

func runCmd(c Command) {
	logF("starting %s: %s", c.Name, c.Command[0])
	cmd := exec.Command(c.Command[0], c.Command[1:]...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	atomic.StoreInt32(keepRunning, 0)
	logF("%s exited", c.Name)
}

func printHelp(tryFiles []string) {
	fmt.Println("starter help")
	fmt.Println("")
	fmt.Println("loads the config, then starts all the configured commands")
	fmt.Println(" if init is specified it runs first, must exit 0")
	fmt.Println(" if any command exits the whole thing exits")
	fmt.Println("")

	fmt.Println("can load config from:")
	for _, p := range tryFiles {
		fmt.Println("", p)
	}
	fmt.Println("")
	c := Config{
		Init: []string{
			"initcommand",
			"initops",
			"can be omitted",
		},
		Commands: []Command{
			{
				Name: "MyCommand",
				Command: []string{
					"command",
					"arg1",
					"arg2",
				},
			},
			{
				Name: "MyCommand2",
				Command: []string{
					"command",
					"arg1",
					"arg2",
				},
			},
		},
	}
	fmt.Println("JSON format")
	jb, _ := json.MarshalIndent(c, "", " ")
	fmt.Println(string(jb))

	fmt.Println("")
	fmt.Println("YAML format")
	yb, _ := yaml.Marshal(c)
	fmt.Println(string(yb))
}

func main() {

	tryFiles := []string{}
	var i0 int32 = 1
	keepRunning = &i0

	// add suffixes to trypaths
	for _, loc := range locations {
		tryFiles = append(tryFiles, loc+".json")
		tryFiles = append(tryFiles, loc+".yaml")
	}

	// help
	if len(os.Args) >= 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			printHelp(tryFiles)
			os.Exit(1)
		}
	}

	// load the config
	config, err := loadConfig(tryFiles)
	if err != nil {
		logF(err.Error())
		os.Exit(1)
	}

	// run init first
	if len(config.Init) != 0 {
		logF("running init: %s", config.Init[0])
		cmd := exec.Command(config.Init[0], config.Init[1:]...)
		out, err := cmd.CombinedOutput()
		fmt.Println(string(out))
		if err != nil {
			logF("init err: %s", err)
			os.Exit(1)
		}
	}

	// if we have no commands just exit
	if len(config.Commands) == 0 {
		logF("no commands to run, exiting")
		os.Exit(0)
	}

	// start all the commands
	for _, c := range config.Commands {
		go runCmd(c)
	}

	for atomic.LoadInt32(keepRunning) == 1 {
		// all still good!
		time.Sleep(time.Millisecond * 500)
	}
	logF("exiting")
}
