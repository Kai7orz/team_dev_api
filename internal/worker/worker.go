package worker

import (
	"log"
	"time"

	"github.com/Kai7orz/team_dev_api/internal/cache"
	"github.com/Kai7orz/team_dev_api/internal/metmuseum"
	"github.com/Kai7orz/team_dev_api/internal/model"
)

//オブジェクトを取得してキャッシュに保存する

func StartWorker(ids []int) {
	go func() {
		for _, id := range ids {

			cache.GlobalCache.Mu.RLock()
			_, exists := cache.GlobalCache.CacheMap[id]
			cache.GlobalCache.Mu.RUnlock()
			if exists {
				continue
			}

			rawArtwork, err := metmuseum.NewClient().GetArtworkByID(id)
			if err != nil || rawArtwork == nil || rawArtwork.ObjectID == 0 {
				log.Print("Failed to fetch artwork id :", id)

				emptyArtwork := model.Artwork{
					ID:           id,
					Title:        nil,
					Artist:       nil,
					Culture:      nil,
					ObjectDate:   nil,
					PrimaryImage: nil,
				}

				//オブジェクト取得できなかった場合は，nilを入れておく
				cache.Save(&emptyArtwork) //取得できないオブジェクトに関してはID情報以外入れない
				time.Sleep(1 * time.Second)
				continue
			}

			artwork := model.Artwork{
				ID:           rawArtwork.ObjectID,
				Title:        &rawArtwork.Title,
				Artist:       &rawArtwork.ArtistDisplayName,
				Culture:      &rawArtwork.Culture,
				ObjectDate:   &rawArtwork.ObjectDate,
				PrimaryImage: &rawArtwork.PrimaryImage,
			}

			cache.Save(&artwork)
			time.Sleep(1 * time.Second)
		}
	}()
}
