package db

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./internal/db/artworks.db")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	var tableExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='artworks')").Scan(&tableExists)
	if err != nil {
		log.Fatalf("failed to check if table exists: %v", err)
	}

	// テーブルが既に存在し、データが入っている場合はインポートをスキップ
	if tableExists {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM artworks").Scan(&count)
		if err != nil {
			log.Fatalf("failed to count records: %v", err)
		}

		if count > 0 {
			log.Printf("Table 'artworks' already exists with %d records. Skipping CSV import.", count)
			return
		}
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
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	csvPath := filepath.Join(cwd, "internal", "metmuseum", "MetObjects.csv")

	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		log.Fatalf("CSVファイルが存在しません: %s", csvPath)
	}

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("failed to open CSV file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close CSV file: %v", err)
		}
	}()

	if _, err := file.Seek(0, 0); err != nil {
		log.Fatalf("failed to seek in CSV file: %v", err)
	}

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	index := map[string]int{}
	for i, h := range headers {
		index[h] = i
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO artworks (
			id,
			accession_year,
			object_begin_date,
			object_end_date,
			artist_begin_date,
			artist_end_date,
			metadata_date,
			is_highlight,
			is_public_domain,
			department,
			object_name,
			title,
			culture,
			period,
			dynasty,
			reign,
			portfolio,
			artist_display_name,
			artist_nationality,
			artist_gender,
			object_date,
			medium,
			classification,
			link_resource,
			artist_wikidata_url,
			object_wikidata_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("failed to close statement: %v", err)
		}
	}()

	count := 0

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		id, _ := strconv.Atoi(row[index["Object ID"]])
		accessionYear, _ := strconv.Atoi(row[index["AccessionYear"]])
		objectBeginDate, _ := strconv.Atoi(row[index["Object Begin Date"]])
		objectEndDate, _ := strconv.Atoi(row[index["Object End Date"]])
		artistBeginDate, _ := strconv.Atoi(row[index["Artist Begin Date"]])
		artistEndDate, _ := strconv.Atoi(row[index["Artist End Date"]])
		metadataDate := row[index["Metadata Date"]]
		isHighlight := strings.ToLower(row[index["Is Highlight"]]) == "true"
		isPublicDomain := strings.ToLower(row[index["Is Public Domain"]]) == "true"

		_, err = stmt.Exec(
			id,
			accessionYear,
			objectBeginDate,
			objectEndDate,
			artistBeginDate,
			artistEndDate,
			metadataDate,
			boolToInt(isHighlight),
			boolToInt(isPublicDomain),
			row[index["Department"]],
			row[index["Object Name"]],
			row[index["Title"]],
			row[index["Culture"]],
			row[index["Period"]],
			row[index["Dynasty"]],
			row[index["Reign"]],
			row[index["Portfolio"]],
			row[index["Artist Display Name"]],
			row[index["Artist Nationality"]],
			row[index["Artist Gender"]],
			row[index["Object Date"]],
			row[index["Medium"]],
			row[index["Classification"]],
			row[index["Link Resource"]],
			row[index["Artist Wikidata URL"]],
			row[index["Object Wikidata URL"]],
		)
		if err != nil {
			log.Printf("Error inserting row with ID %d: %v", id, err)
		}

		count++
		if count%10000 == 0 {
			log.Printf("%d rows processed...", count)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	log.Println("csv import completed")
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
