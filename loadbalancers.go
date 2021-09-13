package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

var LoadBalancerTemplate = `
<table>
<tbody>
<tr>
<th>LoadBalancerName</th>
</tr>
{{range .}}<tr>{{range .}}<td>{{.}}</td>{{end}}</tr>
{{end}}</tbody>
</table>`

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
