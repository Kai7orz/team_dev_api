package model

import (
	"time"
)

func (c *Cache) UpdateLastUsedAt(id int) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.CacheMap[id].LastUsedAt = time.Now().Unix()

}

func (c *Cache) GetCachedDataById(id int) (*Artwork, bool) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	object, ok := c.CacheMap[id]
	if !ok || (time.Now().Unix()-c.CacheMap[id].LastUsedAt) >= int64(c.Ttl) {
		return nil, false
	}
	return object, true
}
