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

	// 生のデータからモデルへの変換
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
