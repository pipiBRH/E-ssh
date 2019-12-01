package sshclient

import (
	"golang.org/x/crypto/ssh"
)

func SSHToIntanceThroughBastion(bastion, target, user string) error {
	bastionClient, err := DialWithSSHAgent(bastion, user)
	if err != nil {
		return err
	}
	defer bastionClient.Close()

	targetClient, err := DialWithJumpHost(bastionClient, target, user)
	if err != nil {
		return err
	}
	defer targetClient.Close()

	if err := targetClient.Terminal(nil).Start(); err != nil {
		return err
	}

	config := &TerminalConfig{
		Term:   "xterm",
		Hight:  40,
		Weight: 80,
		Modes: ssh.TerminalModes{
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		},
	}
	if err := targetClient.Terminal(config).Start(); err != nil {
		return err
	}

	return nil
}
