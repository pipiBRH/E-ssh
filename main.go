package main

import (
	"log"

	"github.com/pipiBRH/E-ssh/cmd"
	"github.com/pipiBRH/E-ssh/sshclient"
	"golang.org/x/crypto/ssh"
)

func main() {
	cmd.Execute()
}

func sshToIntance(bastion, target, user string) {
	bastionClient, err := sshclient.DialWithSSHAgent(bastion, user)
	if err != nil {
		log.Fatal(err)
	}
	defer bastionClient.Close()

	targetClient, err := sshclient.DialWithJumpHost(bastionClient, target, user)
	if err != nil {
		log.Fatal(err)
	}
	defer targetClient.Close()

	if err := targetClient.Terminal(nil).Start(); err != nil {
		log.Fatal(err)
	}

	config := &sshclient.TerminalConfig{
		Term:   "xterm",
		Hight:  40,
		Weight: 80,
		Modes: ssh.TerminalModes{
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		},
	}
	if err := targetClient.Terminal(config).Start(); err != nil {
		log.Fatal(err)
	}
}
