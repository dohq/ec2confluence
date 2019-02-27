package main

import (
	"bytes"
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
	ResultTemplate      bytes.Buffer
	allInstances        [][]string
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

func main() {
	if err := GetInstances(); err != nil {
		log.Fatal(err)
	}

	table, err := RendarTemplate(allInstances)
	if err != nil {
		log.Fatal(err)
	}

	r, err := UpdateContents(confluenceURL, confluenceUSER, confluencePASS, confluencePageTitle, confluencePageID, table)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Update New version %v\n", r)
}

// RendarTemplate is convert instances to Confluence table markup
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

// UpdateContents is login confluence and update wiki page
func UpdateContents(url, user, pass, title, id, table string) (int, error) {
	api, err := goconfluence.NewAPI(url, user, pass)
	if err != nil {
		return 0, errors.Wrap(err, "Cant login confluence")
	}

	c, err := api.GetContentByID(id)
	if err != nil {
		return 0, errors.Wrap(err, "Cant find contents")
	}

	curVersion := c.Version.Number
	newVersion := curVersion + 1

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
			Number: newVersion,
		},
		Space: goconfluence.Space{
			Key: confluencePageSpace,
		},
	}

	content, err := api.UpdateContent(data)
	if err != nil {
		return 0, errors.Wrap(err, "Fail Contents update")
	}

	return content.Version.Number, nil
}
