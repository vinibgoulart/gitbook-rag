package gitbook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Page struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Type        string       `json:"type"`
	Path        string       `json:"path"`
	Document    PageDocument `json:"document"`
}

type PageDocument struct {
	Nodes []map[string]interface{} `json:"nodes"`
}

func SpacesContentPageGet(spaceId *string, pageId *string) (*Page, error) {
	requestURL := GITBOOK_URL + "/spaces/" + *spaceId + "/content" + "/page/" + *pageId

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

	var spaceContentPage Page

	err = json.NewDecoder(res.Body).Decode(&spaceContentPage)
	if err != nil {
		return nil, err
	}

	return &spaceContentPage, nil
}
