package bamboo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	url  string
	user string
	pass string
}

type Result struct {
	Key    string `json:"key"`
	Number int    `json:"buildNumber"`
	State  string `json:"buildState"`

	Plan struct {
		Key       string `json:"key"`
		Name      string `json:"name"`
		ShortKey  string `json:"shortKey"`
		ShortName string `json:"shortName"`
	} `json:"plan"`
}

func New(url, username, password string) *Client {
	return &Client{
		url:  url,
		user: username,
		pass: password,
	}
}

func (c *Client) GetResults(project, plan string) ([]Result, error) {
	url := fmt.Sprintf("%s/rest/api/latest/result/%s-%s.json", c.url, project, plan)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Result{}, err
	}

	req.SetBasicAuth(c.user, c.pass)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []Result{}, err
	}

	results := struct {
		Results struct {
			Results []Result `json:"result"`
		} `json:"results"`
	}{}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return []Result{}, err
	}

	return results.Results.Results, nil
}
