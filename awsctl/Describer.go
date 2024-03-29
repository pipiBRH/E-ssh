package awsctl

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func DescribeEC2(opt *session.Options, input *ec2.DescribeInstancesInput) ([]*ec2.Reservation, error) {
	sess, err := session.NewSessionWithOptions(*opt)
	if err != nil {
		log.Fatal(err)
	}

	svc := ec2.New(sess)

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}
	return result.Reservations, nil
}
