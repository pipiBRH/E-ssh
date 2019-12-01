package interactive

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pipiBRH/E-ssh/sshclient"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pipiBRH/E-ssh/awsctl"
)

func InteractiveWithStdin(profile, region, jumper, user string) {
	option := &session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	if profile != "" {
		option.Profile = profile
	}

	if region != "" {
		option.Config = aws.Config{
			Region: aws.String(region),
		}
	}

	if jumper == "" {
		ec2Des := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name:   aws.String("tag:Name"),
					Values: []*string{aws.String("bastion")},
				},
			},
		}

		ec2Info, err := awsctl.DescribeEC2(option, ec2Des)
		if err != nil {
			log.Fatal(err)
		}

		key, err := shellPrint("Choose a bastion...", "Please enter numeric choice: ", ec2Info)
		if err != nil {
			log.Fatal(err)
		}

		jumper = fmt.Sprintf("%s:2222", *ec2Info[key].Instances[0].PublicIpAddress)
	}

	ec2Des := &ec2.DescribeInstancesInput{}

	ec2Info, err := awsctl.DescribeEC2(option, ec2Des)
	if err != nil {
		log.Fatal(err)
	}

	key, err := shellPrint("Choose a target host...", "Please enter numeric choice: ", ec2Info)
	if err != nil {
		log.Fatal(err)
	}

	target := fmt.Sprintf("%s:22", *ec2Info[key].Instances[0].PrivateIpAddress)

	err = sshclient.SSHToIntanceThroughBastion(jumper, target, user)
	if err != nil {
		log.Fatal(err)
	}
}

func shellPrint(prefix, suffix string, ec2Info []*ec2.Reservation) (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prefix)

	for key, ec2 := range ec2Info {
		ip := *ec2.Instances[0].PublicIpAddress
		dns := *ec2.Instances[0].PublicDnsName
		name := *ec2.Instances[0].Tags[0].Value
		fmt.Printf("(%v) %v, %v, %v\n", key, dns, ip, name)
	}

	fmt.Print(suffix)

	stdin, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	numeric, err := strconv.Atoi(strings.Trim(stdin, "\n"))
	if err != nil {
		return 0, err
	}

	return numeric, nil
}
