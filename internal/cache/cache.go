package cache

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/Kai7orz/team_dev_api/internal/model"
)

var GlobalCache = model.Cache{
	CacheMap: map[int]*model.Artwork{},
	MaxSize:  2,
	Ttl:      60,
}

func GetByID(id int) (*model.Artwork, bool) {

	//キャッシュの中からオブジェクト取得
	GlobalCache.Mu.RLock()

	//Getメソッドでデータがあれば，そのままオブジェクト返して，データの有効期限を更新する
	//キャッシュにあるか見て，無ければDBからとってくる処理を実装

	//データTTL生存チェックしデータ管理させつつ，DBから新しく読み込むときは，キャッシュ内の古いデータから順にリプレースしていく
	object, ok := GlobalCache.GetCachedDataById(id)
	GlobalCache.Mu.RUnlock()
	if !ok {
		singleObject, err := ReadDbByID(id) //キャッシュにないか，キャッシュ内の生存時間過ぎていたら，データベースからデータを取り出す
		if err != nil {
			fmt.Println("Error At Reading Database ")
			return nil, false
		}

		err = Save(singleObject) //キャッシュに保存
		if err != nil {
			fmt.Println("Error At Saving Cache:", err)
		}
		return singleObject, true
	}

	GlobalCache.UpdateLastUsedAt(id) //最終接触時刻の更新

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

func ReadDbByID(id int) (*model.Artwork, error) { //この関数はDBから1つデータ読み込むだけで，キャッシュへの保存はしていないことに注意する

	//以下のCSVファイルから読み込んでいる部分を，DBから読み込むように修正
	cwd, _ := os.Getwd()
	file, err := os.Open(cwd + "/internal/metmuseum/MetObjects.csv")
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	//DBとマッピングされているデータを読み込む（現状は，テストとして1つのデータを定義）
	title := "title"
	artist := "artist"
	culture := "culture"
	objectDate := "objectDate"

	responseObject := &model.Artwork{
		ID:           id,
		LastUsedAt:   time.Now().Unix(),
		Title:        &title,
		Artist:       &artist,
		Culture:      &culture,
		ObjectDate:   &objectDate,
		PrimaryImage: nil,
	}

	return responseObject, nil

}

// キャッシュへの保存処理 ここでキャッシュのデータリプレースなども管理することに注意
func Save(object *model.Artwork) error {

	if object == nil {
		return fmt.Errorf("cannot save nil data")
	}
	GlobalCache.Mu.Lock()
	defer GlobalCache.Mu.Unlock()

	if _, exists := GlobalCache.CacheMap[object.ID]; exists {
		return fmt.Errorf("already saved")
	}

	if len(GlobalCache.CacheMap) >= GlobalCache.MaxSize { //キャッシュ上限数データキャッシュしてる状態で，新たにキャッシュする必要がある場合は，キャッシュ内で最終利用時刻が最小のものをリプレース（この処理はキャッシュサイズ分の計算量を必要とするのでより効率的な処理が欲しい）
		oldestUsedAt := int64(math.MaxInt64)
		deleteId := 0
		for id, obj := range GlobalCache.CacheMap {
			if obj.LastUsedAt < oldestUsedAt {
				oldestUsedAt = obj.LastUsedAt
				deleteId = id
			}
		}
		if deleteId == 0 {
			return fmt.Errorf("error invalid delete id")
		}

		delete(GlobalCache.CacheMap, deleteId)
		GlobalCache.CacheMap[object.ID] = object //キャッシュへの保存
		return nil
	}

	GlobalCache.CacheMap[object.ID] = object

	return nil
}
