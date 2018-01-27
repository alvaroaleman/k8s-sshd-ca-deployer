package main

import (
	"flag"
	//"fmt"
	"io/ioutil"
	"log"
)

func main() {
	var CAURL, CADest, SSHDConfigPath, restartCommand string
	flag.StringVar(&CAURL, "ca-url", "", "The url to fetch the CA cert from")
	flag.StringVar(&CADest, "ca-dest", "", "The patch at which the CA cert should be saved")
	flag.StringVar(&SSHDConfigPath, "sshd-config-path", "/etc/ssh/sshd_config", "The path of your sshd config")
	flag.StringVar(&restartCommand, "restart-command", "systemctl restart sshd", "The command to execute to restart your sshd")
	flag.Parse()

	if CAURL == "" {
		log.Fatal("ca-url parameter must not be empty!")
	} else if CADest == "" {
		log.Fatal("ca-dest parameter must not be empty!")
	} else if SSHDConfigPath == "" {
		log.Fatal("sshd-config-path parameter must not be empty!")
	} else if restartCommand == "" {
		log.Fatal("restart-command parameter must not be empty!")
	}

	sshdConfigRaw, err := ioutil.ReadFile(SSHDConfigPath)
	if err != nil {
		log.Fatalf("Error reading sshd config: '%v'", err)
	}

	sshdConfig := string(sshdConfigRaw)
	log.Println(sshdConfig)
}
