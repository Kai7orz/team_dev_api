package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Kai7orz/team_dev_api/internal/metmuseum"
	"github.com/Kai7orz/team_dev_api/internal/model"
)

// GetArtworkByID godoc
// @Summary Get artwork by ID
// @Description Fetch artwork information from the Met Museum API by ID
// @Tags artworks
// @Produce json
// @Param id path int true "Artwork ID"
// @Success 200 {object} model.Artwork
// @Failure 400 {object} map[string]string
// @Router /artworks/{id} [get]
func GetArtworkByIDHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/artworks/")
	id, err := strconv.Atoi(path)
	if err != nil || path == "" {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	rawArtwork, err := metmuseum.NewClient().GetArtworkByID(id)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch artwork"}`, http.StatusInternalServerError)
		return
	}

	artwork := model.Artwork{
		ID:           rawArtwork.ObjectID,
		Title:        rawArtwork.Title,
		Artist:       rawArtwork.ArtistDisplayName,
		Culture:      rawArtwork.Culture,
		ObjectDate:   rawArtwork.ObjectDate,
		PrimaryImage: rawArtwork.PrimaryImage,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(artwork); err != nil {
		http.Error(w, `{"error":"failed to encode artwork"}`, http.StatusInternalServerError)
	}
}

// GetArtworksHandler godoc
// @Summary Get a list of artworks
// @Description Fetch a paginated list of artworks from the Met API
// @Tags artworks
// @Accept json
// @Produce json
// @Param page query int false "Page number (default is 1)"
// @Success 200 {array} model.Artwork
// @Failure 400 {object} map[string]string
// @Router /artworks [get]
func GetArtworksHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	var artworkPage []model.Artwork //1ページ分の作品を記録するためのスライス
	var page int                    //ページ番号
	pageStr := r.URL.Query().Get("page")

	if pageStr == "" {
		page = 1 //デフォルトで1ページ目の20件を取得
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 { //1ページ目から設定可能とする
			http.Error(w, `{"error":"invalid page number"}`, http.StatusBadRequest)
			return
		}
	}

	artworkPage = metmuseum.NewClient().GetArtworks(page)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(artworkPage); err != nil {
		http.Error(w, `{"error":"failed to encode artwork"}`, http.StatusInternalServerError)
	}
}
