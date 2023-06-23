package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetDatabase(databaseId string) (*Database, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/databases/%s", c.HostURL, databaseId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	database := &Database{}
	err = json.Unmarshal(body, &database)
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (c *Client) CreateDatabase(database Database) (*Database, error) {
	rb, err := json.Marshal(database)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/databases", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newDatabase := Database{}
	err = json.Unmarshal(body, &newDatabase)
	if err != nil {
		return nil, err
	}

	return &newDatabase, nil
}
