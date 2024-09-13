package config

import (
	"log"
	"os/user"
)

var CLI_VERSION string = "1.0.0"

func GetConfigPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir + "/.anvil"
}
