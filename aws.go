package main

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
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
			for _, t := range i.Tags {
				switch *t.Key {
				case "Name":
					Name = *t.Value
				}
			}
			if i.PublicIpAddress == nil {
				i.PublicIpAddress = aws.String("-")
			}
			instance := []string{
				Name,
				*i.InstanceId,
				*i.InstanceType,
				*i.PublicIpAddress,
				*i.PrivateIpAddress,
				*i.Placement.AvailabilityZone,
			}
			allInstances = append(allInstances, instance)
		}
	}
	return nil
}

// GetLoadBalancers is describe lb instances(state Running Only)
func GetLoadBalancers() error {
	svc := elb.New(session.New())
	res, err := svc.DescribeLoadBalancers(nil)
	if err != nil {
		return err
	}

	for _, r := range res.LoadBalancerDescriptions {
		var Instances []string
		var LoadBalancer []string

		LoadBalancerName := *r.LoadBalancerName
		for i := range r.Instances {
			Instances = append(Instances, *r.Instances[i].InstanceId)
		}

		LoadBalancer = append(LoadBalancer, LoadBalancerName, strings.Join(Instances, ", "))
		allLoadbalancers = append(allLoadbalancers, LoadBalancer)
	}
	return nil
}
