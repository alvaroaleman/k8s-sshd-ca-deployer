package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const sshdCaCertParamName = "TrustedUserCAKeys"

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

	CADest, err := filepath.Abs(CADest)
	if err != nil {
		log.Fatalf("Error converting ca-dest to absolute path: '%v'", err)
	}

	SSHDConfigPath, err = filepath.Abs(SSHDConfigPath)
	if err != nil {
		log.Fatalf("Error converting sshd-config-path to absolue path: 'v'", err)
	}

	resp, err := http.Get(CAURL)
	if err != nil {
		log.Fatalf("Error downloading CACert: '%v'", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading certificate body: '%v'", err)
	}

	err = ioutil.WriteFile(CADest, []byte(body), os.FileMode(int(0600)))
	if err != nil {
		log.Fatalf("Error writing CACert: '%v'", err)
	}

	sshdConfigRaw, err := ioutil.ReadFile(SSHDConfigPath)
	if err != nil {
		log.Fatalf("Error reading sshd config: '%v'", err)
	}

	sshdConfig := string(sshdConfigRaw)
	var newSSHDConfig string

	sshdConfigLine := fmt.Sprintf("%s %s", sshdCaCertParamName, CADest)

	if !strings.Contains(sshdConfig, sshdCaCertParamName) {
		newSSHDConfig = fmt.Sprintf("%s\n%s", sshdConfigLine, sshdConfig)
	} else if strings.Contains(sshdConfig, sshdCaCertParamName) && !strings.Contains(sshdConfig, sshdConfigLine) {
		lines := strings.Split(sshdConfig, "\n")
		var finalConfig []string
		finalConfig = append(finalConfig, sshdConfigLine)
		for _, line := range lines {
			if !strings.Contains(line, sshdCaCertParamName) {
				finalConfig = append(finalConfig, line)
			}
		}
		newSSHDConfig = strings.Join(finalConfig, "\n")
	} else if strings.Contains(sshdConfig, sshdConfigLine) {
		newSSHDConfig = sshdConfig
	}

	if newSSHDConfig != sshdConfig {
		err = ioutil.WriteFile(SSHDConfigPath, []byte(newSSHDConfig), os.FileMode(int(0600)))
		if err != nil {
			log.Fatalf("Error writing sshd config: '%v'", err)
		}
		log.Println("Successfully altered sshd config...")

		log.Println("Going to restart sshd....")
		cmdSlice := strings.Split(restartCommand, " ")
		cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Error executing restart command '%s': '%v'\nOutput: %s", restartCommand, err, output)
		}
		log.Println("Successfully reconfigured sshd!")
	} else {
		log.Println("Config already correct, nothing to do...")
	}
}
