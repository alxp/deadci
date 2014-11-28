package main

import (
	"flag"
	"fmt"
	"github.com/dlintw/goconf"
	"log"
	"os"
	"strings"
)

var Config struct {
	DataDir string
	Command []string
	Port    int
	Host    string
	Github  struct {
		Enabled bool
		Token   string
		Secret  string
	}
}

func init() {
	flag.StringVar(&Config.DataDir, "data-dir", "", "Data directory where config is stored. Must be writable.")
}

func InitConfig() {
	flag.Parse()

	// Parse data directory
	if Config.DataDir == "" {
		fmt.Println("  --data-dir option not specified. Please create a writable data directory and invoke deadci command with --data-dir=/path/to/data/dir")
		flag.CommandLine.PrintDefaults()
		os.Exit(2)
	}
	Config.DataDir = strings.TrimRight(Config.DataDir, "/ ")

	// Read the config file
	c, err := goconf.ReadConfigFile(Config.DataDir + "/deadci.ini")
	if err != nil {
		log.Fatal(err.Error() + ". Please ensure that your deadci.ini file is readable and in place at " + Config.DataDir + "/deadci.ini")
	}

	// Parse command
	cmd, err := c.GetString("", "command")
	if err != nil {
		log.Fatal(err)
	}
	cmd = strings.Trim(cmd, " ")
	if cmd == "" {
		log.Fatal("Missing command in deadci.ini. Please specify a command to run to build / test your repositories.")
	}
	Config.Command = strings.Split(cmd, " ")
	if len(Config.Command) == 0 {
		log.Fatal("Missing command in deadci.ini. Please specify a command to run to build / test your repositories.")
	}

	// Parse Port
	Config.Port, err = c.GetInt("", "port")
	if err != nil {
		log.Fatal(err)
	}

	// Parse Host
	Config.Host, err = c.GetString("", "host")
	if (err != nil && err.(goconf.GetError).Reason == goconf.OptionNotFound) || Config.Host == "" {
		Config.Host, err = os.Hostname()
		if err != nil {
			log.Fatal("Unable to determine hostname. Please specify a hostname in deadci.ini")
		}
	} else if err != nil {
		log.Fatal(err)
	}

	// Parse Github settings
	if c.HasSection("github") {
		Config.Github.Enabled, err = c.GetBool("github", "enabled")
		if err != nil && err.(goconf.GetError).Reason != goconf.OptionNotFound {
			log.Fatal(err)
		}
		if Config.Github.Enabled {
			Config.Github.Token, err = c.GetString("github", "token")
			if err != nil && err.(goconf.GetError).Reason != goconf.OptionNotFound {
				log.Fatal(err)
			}
			Config.Github.Secret, err = c.GetString("github", "secret")
			if err != nil && err.(goconf.GetError).Reason != goconf.OptionNotFound {
				log.Fatal(err)
			}
		}
	}

}
