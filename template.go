package main

import (
	"text/template"
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

var LoadBalancerTemplate = `
<table>
<tbody>
<tr>
<th>LoadBalancerName</th>
</tr>
{{range .}}<tr>{{range .}}<td>{{.}}</td>{{end}}</tr>
{{end}}</tbody>
</table>`

// RendarTemplate is convert instances to Confluence table markup
func RendarTemplate(items [][]string, templateText string) (string, error) {
	t, err := template.New("").Parse(templateText)
	if err != nil {
		return "", err
	}

	if err := t.Execute(&ResultTemplate, items); err != nil {
		return "", err
	}

	return ResultTemplate.String(), nil
}
