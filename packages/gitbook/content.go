package gitbook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Content struct {
	ID    string                   `json:"id"`
	Pages []map[string]interface{} `json:"pages"`
}

func ContentGet(spaceId *string) (*Content, error) {
	requestURL := GITBOOK_URL + "/spaces/" + *spaceId + "/content"
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

	var spacesContent Content
	err = json.NewDecoder(res.Body).Decode(&spacesContent)

	if err != nil {
		return nil, err
	}

	return &spacesContent, nil
}
