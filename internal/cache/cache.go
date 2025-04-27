package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/Kai7orz/team_dev_api/internal/metmuseum"
	"github.com/Kai7orz/team_dev_api/internal/model"
)

var GlobalCache = Cache{
	CacheMap: map[int]*model.Artwork{},
	MaxSize:  2,
	Ttl:      60,
}

var GlobalPageCache = PageCache{
	PageData:    make(map[int]Page),
	MaxPageSize: 2,
	Ttl:         6,
}

type Page struct {
	Artworks   []*model.Artwork
	LastUsedAt int64
}

type PageCache struct {
	Mu          sync.RWMutex
	PageData    map[int]Page
	MaxPageSize int
	Ttl         int64
}

type Cache struct {
	Mu       sync.RWMutex
	CacheMap map[int]*model.Artwork
	MaxSize  int
	Ttl      int
}

func GetByID(id int) (*model.Artwork, bool) {

	//キャッシュにあるか見て，無ければDBからとってくる処理を実装
	GlobalCache.Mu.RLock()

	object, ok := GlobalCache.GetCachedDataByID(id) //GetCachedDataByIdメソッドにより，キャッシュ内にデータがあり，生存時間内であればデータが返される
	GlobalCache.Mu.RUnlock()
	if !ok {
		client := metmuseum.NewClient()
		raw, err := client.GetArtworkByID(id) //直接APIたたいてデータを取得する処理
		if err != nil {
			fmt.Printf("error: server internal")
			return nil, false
		}

		responseObject := &model.Artwork{
			ID:           raw.ObjectID,
			LastUsedAt:   time.Now().Unix(),
			Title:        &raw.Title,
			Artist:       &raw.ArtistDisplayName,
			Culture:      &raw.Culture,
			ObjectDate:   &raw.ObjectDate,
			PrimaryImage: &raw.PrimaryImage,
		}

		err = GlobalCache.Save(responseObject) //キャッシュに保存
		if err != nil {
			fmt.Println("Error At Saving Cache:", err)
		}
		return responseObject, true
	}

	GlobalCache.UpdateLastUsedAt(id) //データの有効期限を更新

	return object, true
}
