package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pipiBRH/E-ssh/awsctl"
	"github.com/pipiBRH/E-ssh/cmd"
	"github.com/pipiBRH/E-ssh/sshclient"
	"golang.org/x/crypto/ssh"
)

func main() {
	cmd.Execute()

	if cmd.Jumper == "" {
		option := &session.Options{
			// Specify profile to load for the session's config
			Profile: "aem",

			// Provide SDK Config options, such as Region.
			// Config: aws.Config{
			// 	Region: aws.String("us-west-2"),
			// },

			// Force enable Shared Config support
			SharedConfigState: session.SharedConfigEnable,
		}

		ec2Des := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String("tag:Name"),
					Values: []*string{aws.String("bastion")},
				},
			},
		}

		ec2Info, err := awsctl.Ec2Lookup(option, ec2Des)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
	}

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
