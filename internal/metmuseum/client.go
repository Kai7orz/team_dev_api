// internal/metmuseum/client.go
package metmuseum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Kai7orz/team_dev_api/internal/model"
)

const (
	metMuseumAPIBaseURL = "https://collectionapi.metmuseum.org/public/collection/v1"
)

// RawArtwork はMetMuseum APIから返される生のアートワークデータを表します
type RawArtwork struct {
	ObjectID          int    `json:"objectID"`
	Title             string `json:"title"`
	Culture           string `json:"culture"`
	ArtistDisplayName string `json:"artistDisplayName"`
	ObjectDate        string `json:"objectDate"`
	PrimaryImage      string `json:"primaryImage"`
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

func (c *Client) GetArtworkByID(id int) (*RawArtwork, error) {
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

	var rawArtwork RawArtwork
	if err := json.NewDecoder(resp.Body).Decode(&rawArtwork); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &rawArtwork, nil
}

// limit 個分のデータを返します
func (c *Client) GetArtworks(page int) []model.Artwork {

	var artworks []model.Artwork
	id := 1
	recordedArts := 0 //データの取得が正常に行えた作品数の記録
	limit := 20       //20件のデータを一度のリクエストに対して返す

	for recordedArts < limit {
		rawArtwork, err := NewClient().GetArtworkByID((page-1)*limit + id)
		if err != nil || rawArtwork.ObjectID == 0 {
			id++
			log.Println("error") //情報取得でエラー発生時は，次idのオブジェクトを取得する
			continue
		}

		artwork := model.Artwork{
			ID:           rawArtwork.ObjectID,
			Title:        rawArtwork.Title,
			Artist:       rawArtwork.ArtistDisplayName,
			Culture:      rawArtwork.Culture,
			ObjectDate:   rawArtwork.ObjectDate,
			PrimaryImage: rawArtwork.PrimaryImage,
		}

		artworks = append(artworks, artwork)

		recordedArts++
		id++
	}
	return artworks

}
