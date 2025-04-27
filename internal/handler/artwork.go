package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Kai7orz/team_dev_api/internal/cache"
	"github.com/Kai7orz/team_dev_api/internal/db"
)

const MaxAllowedPage = 1000

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

	if artwork, ok := cache.GetByID(id); ok {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(artwork); err != nil {
			http.Error(w, `{"error":"failed to encode artwork"}`, http.StatusInternalServerError)
		}
		return
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
	var page int //ページ番号
	allowedSortColumns := map[string]bool{
		"title":               true,
		"artist_display_name": true,
	}
	pageStr := r.URL.Query().Get("page")
	filterStr := r.URL.Query().Get("culture")
	sortBy := r.URL.Query().Get("sortBy")

	//デフォルトの設定
	if pageStr == "" {
		page = 1 //デフォルトで1ページ目の20件を取得
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 { //1ページ目から設定可能とする
			http.Error(w, `{"error":"invalid page number"}`, http.StatusBadRequest)
			return
		}
	}

	//不正な値かどうか判定する
	if sortBy != "" {
		if !allowedSortColumns[sortBy] {
			http.Error(w, `{"error":"invalid sort column"}`, http.StatusBadRequest)
			return
		}
	}

	if page > MaxAllowedPage { //ページ指定可能な最大値を設定しておく
		http.Error(w, `{"error":"page too large"}`, http.StatusBadRequest)
		return
	}

	artworks := db.GetArtworks(page, filterStr, sortBy)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(artworks); err != nil {
		http.Error(w, `{"error":"failed to encode artwork"}`, http.StatusInternalServerError)
	}
}
