package gitbook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type GitbookSpace struct {
	Items []GitbookSpaceItems `json:"items"`
}

type GitbookSpaceItems struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type GitbookSpaceContent struct {
	ID    string                     `json:"id"`
	Pages []GitbookSpaceContentPages `json:"pages"`
}

type GitbookSpaceContentPages struct {
	Pages []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"pages"`
}

type GitbookSpaceContentPage struct {
	ID          string                          `json:"id"`
	Title       string                          `json:"title"`
	Description string                          `json:"description"`
	Document    GitbookSpaceContentPageDocument `json:"document"`
}

type GitbookSpaceContentPageDocument struct {
	Nodes []map[string]interface{} `json:"nodes"`
}

func SpacesGet() (*GitbookSpace, error) {
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

	var spaces GitbookSpace
	err = json.NewDecoder(res.Body).Decode(&spaces)
	if err != nil {
		return nil, err
	}

	return &spaces, nil
}

func SpacesContentGet(spaceId *string) (*GitbookSpaceContent, error) {
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

	var spacesContent GitbookSpaceContent
	err = json.NewDecoder(res.Body).Decode(&spacesContent)

	if err != nil {
		return nil, err
	}

	return &spacesContent, nil
}

func SpacesContentPageGet(spaceId *string, pageId *string) (*GitbookSpaceContentPage, error) {
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

	var spaceContentPage GitbookSpaceContentPage

	err = json.NewDecoder(res.Body).Decode(&spaceContentPage)
	if err != nil {
		return nil, err
	}

	return &spaceContentPage, nil
}
