package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
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
	// Classic LoadBalancer
	clb := elb.New(session.New())
	res, err := clb.DescribeLoadBalancers(nil)
	if err != nil {
		return err
	}

	// Application and Network LoadBalancer
	alb := elbv2.New(session.New())
	resv2, err := alb.DescribeLoadBalancers(nil)
	if err != nil {
		return err
	}

	for _, r := range res.LoadBalancerDescriptions {
		var LoadBalancer []string
		var LoadBalancerName string

		LoadBalancerName = *r.LoadBalancerName
		LoadBalancer = append(LoadBalancer, LoadBalancerName)
		allLoadbalancers = append(allLoadbalancers, LoadBalancer)
	}

	for _, r := range resv2.LoadBalancers {
		var LoadBalancer []string
		var LoadBalancerName string

		LoadBalancerName = *r.LoadBalancerName
		LoadBalancer = append(LoadBalancer, LoadBalancerName)
		allLoadbalancers = append(allLoadbalancers, LoadBalancer)
	}
	return nil
}
