package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./internal/db/artworks.db")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS artworks (
		id INTEGER PRIMARY KEY,
		accession_year INTEGER,
		object_begin_date INTEGER,
		object_end_date INTEGER,
		artist_begin_date INTEGER,
		artist_end_date INTEGER,
		metadata_date DATE,
		is_highlight BOOLEAN,
		is_public_domain BOOLEAN,
		department TEXT,
		object_name TEXT,
		title TEXT,
		culture TEXT,
		period TEXT,
		dynasty TEXT,
		reign TEXT,
		portfolio TEXT,
		artist_display_name TEXT,
		artist_nationality TEXT,
		artist_gender TEXT,
		object_date TEXT,
		medium TEXT,
		classification TEXT,
		link_resource TEXT,
		artist_wikidata_url TEXT,
		object_wikidata_url TEXT
	);`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
}
