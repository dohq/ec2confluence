package main

import (
	goconfluence "github.com/cseeger-epages/confluence-go-api"
)

// UpdateContents is login confluence and update wiki page
func UpdateContents(url, user, pass, title, id, table string) (int, error) {
	api, err := goconfluence.NewAPI(url, user, pass)
	if err != nil {
		return 0, err
	}

	c, err := api.GetContentByID(id)
	if err != nil {
		return 0, err
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
		return 0, err
	}

	return content.Version.Number, nil
}
