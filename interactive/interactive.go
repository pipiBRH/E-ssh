package interactive

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
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
		ec2Info, err := awsctl.DescribeEC2(option, &ec2.DescribeInstancesInput{})
		if err != nil {
			log.Fatal(err)
		}

		filterInfo := collectEC2Info(ec2Info, true, "bastion")

		if len(filterInfo) == 0 {
			fmt.Println("Could not found any bastion. Please specific a bastion host...")
			os.Exit(0)
		}

		key, err := shellPrint("Choose a bastion...", "Please enter numeric choice: ", filterInfo)
		if err != nil {
			log.Fatal(err)
		}

		jumper = fmt.Sprintf("%s:2222", filterInfo[key][2])
	}

	ec2Info, err := awsctl.DescribeEC2(option, &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Fatal(err)
	}

	filterInfo := collectEC2Info(ec2Info, false, "")

	key, err := shellPrint("Choose a target host...", "Please enter numeric choice: ", filterInfo)
	if err != nil {
		log.Fatal(err)
	}

	target := fmt.Sprintf("%s:22", filterInfo[key][2])

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

func collectEC2Info(ec2Info []*ec2.Reservation, public bool, filter string) [][]string {
	var result [][]string
	key := 0

	for _, ec2 := range ec2Info {
		var ip, dns string
		valueOfEc2 := reflect.ValueOf(ec2.Instances[0])
		if public {
			fieldOfAddr := valueOfEc2.Elem().FieldByName("PublicIpAddress")
			fieldOfDNS := valueOfEc2.Elem().FieldByName("PublicDnsName")
			if !fieldOfAddr.IsValid() || fieldOfAddr.IsNil() || !fieldOfDNS.IsValid() || fieldOfDNS.IsNil() {
				continue
			}
			ip = *ec2.Instances[0].PublicIpAddress
			dns = *ec2.Instances[0].PublicDnsName
		} else {
			fieldOfAddr := valueOfEc2.Elem().FieldByName("PrivateIpAddress")
			fieldOfDNS := valueOfEc2.Elem().FieldByName("PrivateDnsName")
			if !fieldOfAddr.IsValid() || fieldOfAddr.IsNil() || !fieldOfDNS.IsValid() || fieldOfDNS.IsNil() {
				continue
			}
			ip = *ec2.Instances[0].PrivateIpAddress
			dns = *ec2.Instances[0].PrivateDnsName
		}

		name := findTagsElementByKey(ec2.Instances[0].Tags)
		if filter != "" && strings.Index(name, filter) < 0 {
			continue
		}

		result = append(result, []string{strconv.FormatInt(int64(key), 10), dns, ip, name})
		key++
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
