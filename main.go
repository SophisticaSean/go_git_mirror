package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Repos is a struct representing a JSON configuration repo obj
type Repos struct {
	Name    string `json:"name"`
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
	Repos []Repos `json:"repos"`
}

func main() {
	file, err := ioutil.ReadFile("./conf.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(file))
	var configuration Configuration
	err = json.Unmarshal(file, &configuration)
	if err != nil {
		panic(err)
	}
	fmt.Println(configuration)
	for i := range configuration.Repos {
		fmt.Println(configuration.Repos[i])
	}
}
