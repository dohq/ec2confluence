package main

import (
	"fmt"
	"strings"

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

		LoadBalancerName = aws.StringValue(r.LoadBalancerName)
		LoadBalancer = append(LoadBalancer, LoadBalancerName)
		allLoadbalancers = append(allLoadbalancers, LoadBalancer)
	}

	for _, r := range resv2.LoadBalancers {
		var LoadBalancers []string
		var LoadBalancerName string

		LoadBalancerName = aws.StringValue(r.LoadBalancerName)
		LoadBalancers = append(LoadBalancers, LoadBalancerName)
		allLoadbalancers = append(allLoadbalancers, LoadBalancers)
	}
	return nil
}

func GetSecurityGroup() error {
	svc := ec2.New(session.New())
	res, err := svc.DescribeSecurityGroups(nil)
	if err != nil {
		return err
	}

	for _, e := range res.SecurityGroups {
		var Rule []string
		var IpRanges string
		var Protocol string
		var Port string

		for _, ip := range e.IpPermissions {
			var Ranges []string
			for _, r := range ip.IpRanges {
				Ranges = append(Ranges, aws.StringValue(r.CidrIp))
			}

			Protocol = aws.StringValue(ip.IpProtocol)

			if Protocol == "-1" {
				Protocol = "<strong>Protocol:</strong> All"
			} else {
				Protocol = "<strong>Protocol:</strong> " + Protocol
			}

			Port = fmt.Sprint(aws.Int64Value(ip.ToPort))
			if Port == "0" {
				Port = "<strong>Port:</strong> All"
			} else {
				Port = "<strong>Port:</strong> " + Port
			}

			IpRanges = "<strong>AllowFrom:</strong> <br />" + strings.Join(Ranges, "<br />")

			Rule = append(Rule, Protocol, Port, IpRanges, "<br />---")
		}

		// skip default rule
		if aws.StringValue(e.GroupName) == "default" {
			continue
		}

		SecurityGroups := []string{
			aws.StringValue(e.GroupName),
			aws.StringValue(e.Description),
			strings.Join(Rule, "<br />"),
		}
		allSecurityGroups = append(allSecurityGroups, SecurityGroups)
	}
	return nil
}
