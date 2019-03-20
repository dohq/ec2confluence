package main

import (
	"bytes"
	"flag"
	"log"
	"os"
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
	allLoadbalancers    [][]string
	allSecurityGroups   [][]string
)

func main() {
	var target string
	var inventoryList [][]string
	var useTemplate string

	flag.StringVar(&target, "t", "", "Export inventory target(ec2 or lb)")
	flag.Parse()

	switch target {
	case "ec2":
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

	default:
		log.Println("Please set t Args")
		os.Exit(1)
	}

	Table, err := RendarTemplate(inventoryList, useTemplate)
	if err != nil {
		log.Fatalf("Template Error: %v", err)
	}

	r, err := UpdateContents(confluenceURL, confluenceUSER, confluencePASS, confluencePageTitle, confluencePageID, Table)
	if err != nil {
		log.Fatalf("UpdateContents Error: %v", err)
	}

	log.Printf("Update New version %v\n", r)
}
