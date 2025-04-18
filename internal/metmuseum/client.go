// internal/metmuseum/client.go
package metmuseum

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	metMuseumAPIBaseURL = "https://collectionapi.metmuseum.org/public/collection/v1"
)

type Artwork struct {
	ObjectID   int    `json:"objectID"`
	Title      string `json:"title"`
	Culture    string `json:"culture"`
	ArtistName string `json:"artistDisplayName"`
	// 他にも必要なフィールドがあれば追加
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL:    metMuseumAPIBaseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetArtworkByID(id int) (*Artwork, error) {
	url := fmt.Sprintf("%s/objects/%d", c.BaseURL, id)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("error closing response body: %v\n", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	var artwork Artwork
	if err := json.NewDecoder(resp.Body).Decode(&artwork); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &artwork, nil
}
