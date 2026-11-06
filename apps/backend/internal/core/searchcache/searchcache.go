package searchcache

import (
	"encoding/json"
	"time"

	"tunetrend-backend/internal/core/upstash"
	"tunetrend-backend/internal/domain"
)

const (
	ResultTTL   = 24 * time.Hour
	HitDedupTTL = time.Hour
)

type Cache struct {
	store *upstash.Client
}

func New(store *upstash.Client) *Cache {
	return &Cache{store: store}
}

func (c *Cache) Get(key string) ([]domain.SearchItem, bool, error) {
	raw, ok, err := c.store.Get(key)
	if err != nil || !ok {
		return nil, false, err
	}
	var items []domain.SearchItem
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil, false, nil
	}
	return items, true, nil
}

func (c *Cache) Set(key string, items []domain.SearchItem) error {
	raw, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return c.store.SetEX(key, string(raw), ResultTTL)
}

type HitDeduper struct {
	store *upstash.Client
}

func NewHitDeduper(store *upstash.Client) *HitDeduper {
	return &HitDeduper{store: store}
}

func (d *HitDeduper) FirstHit(ip, videoID string) (bool, error) {
	return d.store.SetNX("search:hit:"+ip+":"+videoID, "1", HitDedupTTL)
}
