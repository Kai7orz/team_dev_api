package cache

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/Kai7orz/team_dev_api/internal/model"
)

var GlobalCache = model.Cache{
	CacheMap: map[int]*model.Artwork{},
}

func GetByID(id int) (*model.Artwork, bool) {
	//キャッシュの中からオブジェクト取得
	GlobalCache.Mu.RLock()
	defer GlobalCache.Mu.RUnlock()

	object, ok := GlobalCache.CacheMap[id]
	if !ok {
		return nil, false
	}
	return object, true
}

func GetByPage(page int) []*model.Artwork {
	//キャッシュの中からオブジェクト取得（ページ単位）
	limit := 20
	start := (page-1)*limit + 1
	end := start + limit - 1

	GlobalCache.Mu.RLock()
	defer GlobalCache.Mu.RUnlock()

	var result []*model.Artwork

	for id := start; id <= end; id++ {
		art, ok := GlobalCache.CacheMap[id]
		if ok {
			result = append(result, art)
		} else {
			result = append(result, nil) // キャッシュに存在しないデータはnil
		}
	}
	return result
}

// 検索機能などにおいてnil入ったデータを扱いたくないので，それ等除外したデータを返す
func GetAll() []*model.Artwork {

	GlobalCache.Mu.RLock()
	defer GlobalCache.Mu.RUnlock()

	var result []*model.Artwork

	for _, art := range GlobalCache.CacheMap {
		if art != nil && art.Title != nil {
			result = append(result, art)
		}
	}

	return result
}

func ReadCsv() {
	cwd, _ := os.Getwd()
	file, err := os.Open(cwd + "/internal/metmuseum/MetObjects.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV data:", err)
		return
	}

	limit := 1000 //1000件分のデータを登録
	id := 0

	//以下　今回は手動で各カラムが何番の列に対応するかを数えたが，自動化したほうがいい
	strObjectID := 4
	title := 9
	artist := 18
	culture := 10
	objectDate := 28

	for _, record := range records {
		if id > limit {
			break
		}
		if id == 0 { //csv1行目は破棄
			id++
			continue
		}

		objectID, _ := strconv.Atoi(record[strObjectID])

		tempObject := &model.Artwork{
			ID:           objectID,
			Title:        &record[title],
			Artist:       &record[artist],
			Culture:      &record[culture],
			ObjectDate:   &record[objectDate],
			PrimaryImage: nil,
		}

		Save(tempObject)
		id++
	}
}

// キャッシュへの保存処理
func Save(object *model.Artwork) {
	if object == nil {
		return
	}

	GlobalCache.Mu.Lock()
	defer GlobalCache.Mu.Unlock()

	if _, exists := GlobalCache.CacheMap[object.ID]; exists {
		return
	}

	GlobalCache.CacheMap[object.ID] = object
}
