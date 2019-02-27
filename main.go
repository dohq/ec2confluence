package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	goconfluence "github.com/cseeger-epages/confluence-go-api"
	"github.com/pkg/errors"
)

var (
	confluenceURL       string = os.Getenv("CONFLUENCE_URL")
	confluenceUSER      string = os.Getenv("CONFLUENCE_USER")
	confluencePASS      string = os.Getenv("CONFLUENCE_PASS")
	confluencePageID    string = os.Getenv("CONFLUENCE_PAGE_ID")
	confluencePageTitle string = os.Getenv("CONFLUENCE_PAGE_TITLE")
	confluencePageSpace string = os.Getenv("CONFLUENCE_PAGE_SPACE")
)

var InstancesTemplate = `
<table>
<tbody>
<tr>
<th>InstanceName</th>
<th>InstanceId</th>
<th>InstanceType</th>
<th>PublicIpAddress</th>
<th>PrivateIpAddress</th>
<th>AvailabilityZone</th>
</tr>
{{range .}}<tr>{{range .}}<td>{{.}}</td>{{end}}</tr>
{{end}}</tbody>
</table>`

var (
	ResultTemplate bytes.Buffer
	allInstances   [][]string
)

func main() {

	if err := getInstances(); err != nil {
		log.Fatal(err)
	}

	table, err := RendarTemplate(allInstances)
	if err != nil {
		log.Fatal(err)
	}

	api, err := goconfluence.NewAPI(
		confluenceURL,
		confluenceUSER,
		confluencePASS,
	)
	if err != nil {
		log.Fatal(err)
	}

	c, err := api.GetContentByID(confluencePageID)
	if err != nil {
		log.Fatal(err)
	}

	v := c.Version.Number

	data := &goconfluence.Content{
		ID:    confluencePageID,
		Type:  "page",
		Title: confluencePageTitle,

		Body: goconfluence.Body{
			Storage: goconfluence.Storage{
				Value:          table,
				Representation: "storage",
			},
		},
		Version: goconfluence.Version{
			Number: v + 1,
		},
		Space: goconfluence.Space{
			Key: confluencePageSpace,
		},
	}

	content, err := api.UpdateContent(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(content.Version.Number)
}

// RendarTemplate is convert instances json to Confluence table markup
func RendarTemplate(instances [][]string) (string, error) {
	t, err := template.New("").Parse(InstancesTemplate)
	if err != nil {
		return "", errors.Wrap(err, "Can't Parse Template")
	}

	if err := t.Execute(&ResultTemplate, &instances); err != nil {
		return "", errors.Wrap(err, "Can't Parse Template")
	}

	return ResultTemplate.String(), nil
}

func getInstances() error {
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
		return errors.Wrap(err, "Cant describe instances")
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
				i.PublicIpAddress = aws.String("Not Assigned")
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
