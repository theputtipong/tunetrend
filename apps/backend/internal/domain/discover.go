package domain

type DiscoverItem struct {
	CategoryID    string `json:"categoryId"`
	CategoryLabel string `json:"categoryLabel"`
	ID            string `json:"id"`
	Title         string `json:"title"`
	ChannelTitle  string `json:"channelTitle"`
	ThumbnailURL  string `json:"thumbnailUrl"`
	ViewCount     string `json:"viewCount"`
	CountryCode   string `json:"countryCode"`
}

type DiscoverUsecase interface {
	GetDiscoverItems() ([]DiscoverItem, error)
}
