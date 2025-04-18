// @title           Team Dev API
// @version         1.0
// @description     metmuseumのデータを取得するAPI
// @host            localhost:8080
// @BasePath        /
package main

import (
	"fmt"
	"net/http"

	_ "github.com/Kai7orz/team_dev_api/docs"
	"github.com/Kai7orz/team_dev_api/internal/handler"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	http.Handle("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/artworks/", handler.GetArtworkByIDHandler)

	fmt.Println("Server is running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}
