package gitbook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Space struct {
	Items []SpaceItems `json:"items"`
}

type SpaceItems struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func SpacesGet() (*Space, error) {
	requestURL := GITBOOK_URL + "/orgs/" + os.Getenv("GITBOOK_ORGANIZATION_ID") + "/spaces"
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITBOOK_API_TOKEN"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error: %d", res.StatusCode)
	}

	var spaces Space
	err = json.NewDecoder(res.Body).Decode(&spaces)
	if err != nil {
		return nil, err
	}

	return &spaces, nil
}
