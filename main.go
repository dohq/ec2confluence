package main

import (
	"bytes"
	"log"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	loarBalancers     = kingpin.Command("lb", "Output LoadBalancers Info")
	instances         = kingpin.Command("in", "Output Instances Info")
	secrityGroups     = kingpin.Command("sg", "Output SecurityGroups Info")
	URL               = kingpin.Flag("url", "Confluence URL").Envar("CONFLUENCE_URL").Short('u').Required().String()
	USERNAME          = kingpin.Flag("username", "Confluence UserName").Envar("CONFLUENCE_USER").Short('n').Required().String()
	PASSWORD          = kingpin.Flag("password", "Confluence Password").Envar("CONFLUENCE_PASSWORD").Short('p').Required().String()
	PageID            = kingpin.Flag("id", "Confluence PageID").Envar("CONFLUENCE_PAGE_ID").Short('i').Required().String()
	PageTitle         = kingpin.Flag("title", "Confluence PageTitle").Envar("CONFLUENCE_PAGE_TITLE").Short('t').Required().String()
	PageSpace         = kingpin.Flag("space", "Confluence Space").Envar("CONFLUENCE_PAGE_SPACE").Short('s').Required().String()
	ResultTemplate    bytes.Buffer
	allInstances      [][]string
	allLoadbalancers  [][]string
	allSecurityGroups [][]string
)

func main() {
	var inventoryList [][]string
	var useTemplate string

	switch kingpin.Parse() {
	case "in":
		if err := GetInstances(); err != nil {
			log.Fatalf("Instances Error: %v", err)
		}
		inventoryList = allInstances
		useTemplate = InstancesTemplate

	case "lb":
		if err := GetLoadBalancers(); err != nil {
			log.Fatalf("LoadBalancer Error: %v", err)
		}
		inventoryList = allLoadbalancers
		useTemplate = LoadBalancerTemplate

	case "sg":
		if err := GetSecurityGroup(); err != nil {
			log.Fatalf("SecurityGroup Error: %v", err)
		}
		inventoryList = allSecurityGroups
		useTemplate = SecurityGroupTemplate
	}

	Contents, err := RendarTemplate(inventoryList, useTemplate)
	if err != nil {
		log.Fatalf("Rendaring Template Error: %v", err)
	}

	r, err := UpdateContents(*URL, *USERNAME, *PASSWORD, *PageTitle, *PageID, Contents)
	if err != nil {
		log.Fatalf("UpdateContents Error: %v", err)
	}

	log.Printf("Update New version %v\n", r)
}
