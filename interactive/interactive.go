package interactive

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
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

		filteInfo := collectEC2Info(ec2Info, true)

		key, err := shellPrint("Choose a bastion...", "Please enter numeric choice: ", filteInfo)
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

	filteInfo := collectEC2Info(ec2Info, false)

	key, err := shellPrint("Choose a target host...", "Please enter numeric choice: ", filteInfo)
	if err != nil {
		log.Fatal(err)
	}

	target := fmt.Sprintf("%s:22", *ec2Info[key].Instances[0].PrivateIpAddress)

	err = sshclient.SSHToIntanceThroughBastion(jumper, target, user)
	if err != nil {
		log.Fatal(err)
	}
}

func shellPrint(prefix, suffix string, ec2Info [][]string) (int, error) {
	fmt.Println(prefix)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Dns", "IP", "Name"})

	for _, v := range ec2Info {
		table.Append(v)
	}
	table.Render()

	fmt.Print(suffix)

	reader := bufio.NewReader(os.Stdin)
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

func collectEC2Info(ec2Info []*ec2.Reservation, public bool) [][]string {
	var result [][]string
	for key, ec2 := range ec2Info {
		var ip, dns string
		if public {
			ip = *ec2.Instances[0].PrivateIpAddress
			dns = *ec2.Instances[0].PrivateDnsName

		} else {
			ip = *ec2.Instances[0].PrivateIpAddress
			dns = *ec2.Instances[0].PrivateDnsName
		}
		name := findTagsElementByKey(ec2.Instances[0].Tags)
		result = append(result, []string{strconv.FormatInt(int64(key), 10), dns, ip, name})
	}
	return result
}

func findTagsElementByKey(tags []*ec2.Tag) string {
	var name string
	for _, t := range tags {
		if *t.Key == "Name" {
			name = *t.Value
			break
		}
	}
	return name
}
