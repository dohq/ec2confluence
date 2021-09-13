package main

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// GetInstances is describe ec2 instances(state Running Only)
func GetInstances() error {
	svc := ec2.New(session.New())
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

	res, err := svc.DescribeInstances(params)
	if err != nil {
		return err
	}

	for _, r := range res.Reservations {
		for _, i := range r.Instances {
			var Name string
			var SecurityGroups []string

			for _, t := range i.Tags {
				switch aws.StringValue(t.Key) {
				case "Name":
					Name = aws.StringValue(t.Value)
				}
			}
			if i.PublicIpAddress == nil {
				i.PublicIpAddress = aws.String("-")
			}

			for _, s := range i.SecurityGroups {
				SecurityGroups = append(SecurityGroups, aws.StringValue(s.GroupName))
			}

			instance := []string{
				Name,
				aws.StringValue(i.InstanceId),
				aws.StringValue(i.InstanceType),
				aws.StringValue(i.PublicIpAddress),
				aws.StringValue(i.PrivateIpAddress),
				aws.StringValue(i.Placement.AvailabilityZone),
				strings.Join(SecurityGroups, " / "),
			}
			allInstances = append(allInstances, instance)
		}
	}
	return nil
}
