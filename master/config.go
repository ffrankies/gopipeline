package master

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config stores the pipeline configuration, read in from a YAML file
type Config struct {
	SSHUser  string   `yaml:"SSHUser"`  // The User with which to log into pipeline worker nodes
	SSHPort  int      `yaml:"SSHPort"`  // The port number with which to log into pipelined worker nodes
	NodeList []string `yaml:"NodeList"` // The list of nodes available to the pipeline
	UserPath string   `yaml:"UserPath"` // The userpath for the go install directory
}

// NewConfig creates a new Config object out of a YAMl config file
func NewConfig(configPath string) *Config {
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	config := Config{}
	err = yaml.Unmarshal([]byte(configData), &config)
	if err != nil {
		panic("Could not parse config file")
	}
	return &config
}
