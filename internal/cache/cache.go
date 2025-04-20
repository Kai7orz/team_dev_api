package cache

import (
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
