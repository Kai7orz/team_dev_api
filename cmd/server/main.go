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
	"github.com/Kai7orz/team_dev_api/internal/db"
	"github.com/Kai7orz/team_dev_api/internal/handler"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	db.InitDB()
	http.Handle("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/artworks/", handler.GetArtworkByIDHandler)
	http.HandleFunc("/artworks", handler.GetArtworksHandler)

	fs := http.FileServer(http.Dir("./demo"))
	http.Handle("/demo/", http.StripPrefix("/demo/", fs))

	http.HandleFunc("/", staticHTMLHandler)

	fmt.Println("Server is running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func staticHTMLHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./demo/index.html")
}
