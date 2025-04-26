package cache

import (
	"fmt"
	"math"
	"time"

	"github.com/Kai7orz/team_dev_api/internal/metmuseum"
	"github.com/Kai7orz/team_dev_api/internal/model"
)

var GlobalCache = model.Cache{
	CacheMap: map[int]*model.Artwork{},
	MaxSize:  2,
	Ttl:      60,
}

var GlobalPageCache = model.PageCache{
	PageData:    make(map[int]model.Page),
	MaxPageSize: 2,
	Ttl:         6,
}

func GetByID(id int) (*model.Artwork, bool) {

	//キャッシュにあるか見て，無ければDBからとってくる処理を実装
	GlobalCache.Mu.RLock()

	object, ok := GlobalCache.GetCachedDataByID(id) //GetCachedDataByIdメソッドにより，キャッシュ内にデータがあり，生存時間内であれば，そのままオブジェクト返してデータの有効期限を更新する
	GlobalCache.Mu.RUnlock()
	if !ok {
		singleObject, err := GetDataByID(id) //キャッシュにないか，キャッシュ内の生存時間過ぎていたら，データベースからデータを取り出すーー＞APIからとってきてキャッシュする
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

func GetDataByID(id int) (*model.Artwork, error) { //この関数は外部APIから1つデータ読み込んで返す関数

	client := metmuseum.NewClient()
	raw, err := client.GetArtworkByID(id)
	if err != nil {
		fmt.Printf("error: server internal")
	}

	responseObject := &model.Artwork{
		ID:           raw.ObjectID,
		LastUsedAt:   time.Now().Unix(),
		Title:        &raw.Title,
		Artist:       &raw.ArtistDisplayName,
		Culture:      &raw.ArtistDisplayName,
		ObjectDate:   &raw.PrimaryImage,
		PrimaryImage: &raw.PrimaryImage,
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

	//キャッシュリフレッシュ処理（キャッシュ内の生存時間が過ぎたデータをすべて削除する）
	GlobalCache.Refresh()

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
