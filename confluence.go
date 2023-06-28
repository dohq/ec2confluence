package main

import (
	goconfluence "github.com/virtomize/confluence-go-api"
)

// UpdateContents is login confluence and update wiki page
func UpdateContents(url, user, pass, title, id, space, table string) (int, error) {
	api, err := goconfluence.NewAPI(url, user, pass)
	if err != nil {
		return 0, err
	}

	c, err := api.GetContentByID(id, goconfluence.ContentQuery{})
	if err != nil {
		return 0, err
	}

	curVersion := c.Version.Number
	newVersion := curVersion + 1

	data := &goconfluence.Content{
		ID:    id,
		Type:  "page",
		Title: title,

		Body: goconfluence.Body{
			Storage: goconfluence.Storage{
				Value:          table,
				Representation: "storage",
			},
		},
		Version: &goconfluence.Version{
			Number: newVersion,
		},
		Space: &goconfluence.Space{
			Key: space,
		},
	}

	content, err := api.UpdateContent(data)
	if err != nil {
		return 0, err
	}

	return content.Version.Number, nil
}
