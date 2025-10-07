package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var SecurityGroupTemplate = `
<table>
<tbody>
<tr>
<th>GroupName</th>
<th>Description</th>
<th>Rule</th>
</tr>
{{range .}}<tr>{{range .}}<td>{{.}}</td>{{end}}</tr>
{{end}}</tbody>
</table>`

func GetSecurityGroup() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	svc := ec2.New(sess)
	res, err := svc.DescribeSecurityGroups(nil)
	if err != nil {
		return err
	}

	for _, e := range res.SecurityGroups {
		var Rule []string
		var IPRanges string
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

			IPRanges = "<strong>AllowFrom:</strong> <br />" + strings.Join(Ranges, "<br />")

			Rule = append(Rule, Protocol, Port, IPRanges, "<br />---")
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
