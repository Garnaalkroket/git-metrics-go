package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var ENV_VARIABLE_MAP map[string]string

func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func CheckIfError(err error) {
	if err == nil {
		return
	}
	if err.Error() == "already up-to-date" {
		return
	}
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func InitEnvVariable(variable_name string) {
	env_temp := os.Getenv(variable_name)
	if env_temp == "" {
		Info(variable_name + " not overriden, using default value: " + ENV_VARIABLE_MAP[variable_name])
	} else {
		env_default := ENV_VARIABLE_MAP[variable_name]
		ENV_VARIABLE_MAP[variable_name] = env_temp
		Info(variable_name + " overriden, using value: " + ENV_VARIABLE_MAP[variable_name] + " instead of default " + env_default)
	}
}

func InitEnvVariables() {
	for key := range ENV_VARIABLE_MAP {
		InitEnvVariable(key)
	}
}

func main() {
	ENV_VARIABLE_MAP = make(map[string]string)
	ENV_VARIABLE_MAP["GIT_CLONE_DIR"] = "/tmp/clone-dir"
	ENV_VARIABLE_MAP["GIT_REMOTE_REPO"] = "github.com:Garnaalkroket/git-metrics-go.git"
	ENV_VARIABLE_MAP["GIT_PULL_DIR"] = "/tmp/pull-dir"
	ENV_VARIABLE_MAP["GIT_TIMEOUT_SECONDS"] = "60"
	ENV_VARIABLE_MAP["GIT_SSH_KEY"] = os.Getenv("HOME") + "/.ssh/id_rsa"
	ENV_VARIABLE_MAP["GIT_ACTIONS"] = "clone,pull"
	ENV_VARIABLE_MAP["LOG_FILE_LOCATION"] = "/tmp/tst.txt"
	InitEnvVariables()
	actions := strings.Split(ENV_VARIABLE_MAP["GIT_ACTIONS"], ",")
	var metrics []Measurer
	key, err := ssh.NewPublicKeysFromFile("git", ENV_VARIABLE_MAP["GIT_SSH_KEY"], "")
	if err != nil {
		panic(err)
	}
	for _, action := range actions {

		if action == "pull" {
			metrics = append(metrics, PullAction{
				Action: Action{
					method:   "PULL",
					repo:     ENV_VARIABLE_MAP["GIT_REMOTE_REPO"],
					dir:      ENV_VARIABLE_MAP["GIT_PULL_DIR"],
					log_file: ENV_VARIABLE_MAP["LOG_FILE_LOCATION"],
				},
				key: key,
			})
		}

		if action == "clone" {
			metrics = append(metrics, CloneAction{
				Action: Action{
					method:   "CLONE",
					repo:     ENV_VARIABLE_MAP["GIT_REMOTE_REPO"],
					dir:      ENV_VARIABLE_MAP["GIT_CLONE_DIR"],
					log_file: ENV_VARIABLE_MAP["LOG_FILE_LOCATION"],
				},
				key: key,
			})
		}

	}
	for _, action := range metrics {
		action.InitMeasurement()
	}
	for 1 == 1 {
		for _, action := range metrics {
			action.Measure()
		}
	}
}
