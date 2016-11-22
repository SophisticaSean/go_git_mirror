package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Repos is a struct representing a JSON configuration repo obj
type Repos struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Source  Repo   `json:"source"`
	Mirrors []Repo `json:"mirrors"`
}

// Repo represents a single git repository
type Repo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Configuration is the root struct our configuration file
type Configuration struct {
	Repos    []Repos `json:"repos"`
	Interval int     `json:"interval"`
	HomePath string  `json:"pathToThisFile"`
}

var configuration Configuration

func blahinit() {
	file, err := ioutil.ReadFile("./conf.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &configuration)
	if err != nil {
		panic(err)
	}
}

func updateConfiguration(homePath string) {
	errOne := os.Chdir(homePath)
	if errOne != nil {
		panic(errOne)
	}
	file, err := ioutil.ReadFile("./conf.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &configuration)
	if err != nil {
		panic(err)
	}
}

func commandWrapper(firstArg string, commands []string) (out string, newErr string) {
	cmd := exec.Command(firstArg, commands...)
	var buf bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &stdErr

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	cmd.Wait()
	out = buf.String()
	newErr = stdErr.String()
	return out, newErr
}

func mirrorInit(repos []Repos) {
	for i := range configuration.Repos {
		repo := configuration.Repos[i]
		os.Chdir(repo.Path)
		commandWrapper("git", []string{"remote", "add", repo.Source.Name, repo.Source.URL})
		for i := range repo.Mirrors {
			mirror := repo.Mirrors[i]
			commandWrapper("git", []string{"remote", "add", mirror.Name, mirror.URL})
		}
	}
}

func main() {
	fmt.Println("running")
	blahinit()
	mirrorInit(configuration.Repos)
	for {
		updateConfiguration(configuration.HomePath)
		for i := range configuration.Repos {
			repo := configuration.Repos[i]
			mirrors := repo.Mirrors
			err := os.Chdir(repo.Path)
			if err != nil {
				panic(err)
			}
			output, _ := commandWrapper("git", []string{"pull", repo.Source.Name, "master"})
			if !strings.Contains(output, "up-to-date") {
				fmt.Println("Picked up changes from " + repo.Source.URL)
				fmt.Println(output)
			}
			for i := range mirrors {
				mirror := mirrors[i]
				output, err := commandWrapper("git", []string{"push", mirror.Name, "master"})
				if !strings.Contains(err, "up-to-date") {
					fmt.Println("Pushing changes out to " + mirror.Name)
					fmt.Println(output, err)
				}
			}
			//fmt.Println(configuration.Repos[i])
		}
		time.Sleep(time.Duration(configuration.Interval) * time.Second)
	}
}
