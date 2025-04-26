package model

import (
	"fmt"
	"time"
)

func (c *Cache) UpdateLastUsedAt(id int) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.CacheMap[id].LastUsedAt = time.Now().Unix()
}

func (c *Cache) GetCachedDataByID(id int) (*Artwork, bool) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	object, ok := c.CacheMap[id]
	if !ok || (time.Now().Unix()-c.CacheMap[id].LastUsedAt) >= int64(c.Ttl) {
		return nil, false
	}
	return object, true
}

func (c *Cache) Refresh() {

	for key, value := range c.CacheMap {
		if (time.Now().Unix() - value.LastUsedAt) >= int64(c.Ttl) {
			delete(c.CacheMap, key)
		}
	}
}

func TestLoader(page int) []*Artwork {
	var artworks []*Artwork

	startID := (page-1)*20 + 1

	for i := 0; i < 20; i++ {
		id := startID + i
		title := fmt.Sprintf("Test Title %d", id)
		artist := fmt.Sprintf("Test Artist %d", id)
		culture := fmt.Sprintf("Test Culture %d", id)
		objectDate := "2025"
		primaryImage := fmt.Sprintf("https://example.com/image%d.jpg", id)

		artworks = append(artworks, &Artwork{
			ID:           id,
			LastUsedAt:   time.Now().Unix(),
			Title:        &title,
			Artist:       &artist,
			Culture:      &culture,
			ObjectDate:   &objectDate,
			PrimaryImage: &primaryImage,
		})
	}

	return artworks
}

func (pc *PageCache) GetPage(page int) []*Artwork {
	//はじめにキャッシュからデータを探して，無ければDBからとってくる
	pc.Mu.Lock()
	defer pc.Mu.Unlock()

	now := time.Now().Unix()

	p, ok := pc.PageData[page]
	if !ok || len(p.Artworks) == 0 || now-p.LastUsedAt > pc.Ttl {
		fmt.Println("[Cache Miss] page:", page)

		artworks := TestLoader(page) //DBからデータを取得
		if artworks != nil {
			pc.savePageInternal(page, artworks, now)
		}
		return artworks
	}
	fmt.Println("[Cache Hit] page:", page)

	pc.PageData[page] = Page{
		Artworks:   p.Artworks,
		LastUsedAt: now,
	}

	return p.Artworks

}

func (pc *PageCache) DeleteOldestPage() {
	oldestPage := -1
	oldestTime := time.Now().Unix()

	for pageNum, p := range pc.PageData {

		if len(p.Artworks) > 0 && p.LastUsedAt <= oldestTime {
			oldestTime = p.LastUsedAt
			oldestPage = pageNum
		}
	}

	if oldestPage != -1 {
		delete(pc.PageData, oldestPage)
	}

}

func (pc *PageCache) CountNonEmptyPages() int {
	count := 0
	for _, p := range pc.PageData {
		if len(p.Artworks) > 0 { //中身があるか判定
			count++
		}
	}
	return count
}

func (pc *PageCache) savePageInternal(page int, artworks []*Artwork, now int64) {

	if len(pc.PageData) >= pc.MaxPageSize {
		pc.DeleteOldestPage()
	}

	pc.PageData[page] = Page{
		Artworks:   artworks,
		LastUsedAt: now,
	}

}

func (pc *PageCache) SavePage(page int, artworks []*Artwork) {
	pc.Mu.Lock()
	defer pc.Mu.Unlock()

	now := time.Now().Unix()

	//PageCacheのPageData[i]のiがページ番号に対応しているので，page が既存のiより大きければ，expandによってキャッシュ用のメモリ領域を新たに確保する必要がある

	if pc.CountNonEmptyPages() >= pc.MaxPageSize { //キャッシュ容量が上限超えていれば，古いデータを新データでリプレース
		pc.DeleteOldestPage()
		return
	}

	pc.PageData[page] = Page{
		Artworks:   artworks,
		LastUsedAt: now,
	}

}
