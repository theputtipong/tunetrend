package usecase

import (
	"fmt"
	"log"

	"tunetrend-backend/internal/domain"
)

type songUsecase struct {
	songRepo      domain.SongRepository
	ytRepo        domain.YouTubeRepository
	categoryRepos map[string]domain.CategoryVideoRepository
}

func NewSongUsecase(songRepo domain.SongRepository, ytRepo domain.YouTubeRepository, categoryRepos map[string]domain.CategoryVideoRepository) domain.SongUsecase {
	return &songUsecase{
		songRepo:      songRepo,
		ytRepo:        ytRepo,
		categoryRepos: categoryRepos,
	}
}

func (u *songUsecase) SyncTrendingMusic(countryCode string) error {
	log.Printf("🔄 [Usecase] กำลังดึงข้อมูล YouTube เทรนด์สำหรับประเทศ %s...", countryCode)

	songs, err := u.ytRepo.FetchTrending(countryCode)
	if err != nil {
		return err
	}

	err = u.songRepo.UpsertSongs(songs)
	if err != nil {
		return err
	}

	log.Printf("✅ [Usecase] ซิงก์ข้อมูล %s สำเร็จ จำนวน %d เพลง", countryCode, len(songs))
	return nil
}

func (u *songUsecase) resolveCategoryRepo(categoryID string) (domain.CategoryVideoRepository, error) {
	catRepo, ok := u.categoryRepos[categoryID]
	if !ok {
		return nil, fmt.Errorf("%w: %s", domain.ErrUnknownCategory, categoryID)
	}
	return catRepo, nil
}

func mapCategoryVideosToSongs(videos []domain.CategoryVideo, categoryID string) []domain.Song {
	songs := make([]domain.Song, len(videos))
	for i, v := range videos {
		songs[i] = domain.Song{
			ID:           v.ID,
			Title:        v.Title,
			ChannelTitle: v.ChannelTitle,
			ThumbnailURL: v.ThumbnailURL,
			ViewCount:    v.ViewCount,
			CountryCode:  v.CountryCode,
			CategoryID:   categoryID,
			PublishedAt:  v.PublishedAt,
		}
	}
	return songs
}

func (u *songUsecase) GetTrends(countryCode, categoryID string) ([]domain.Song, error) {
	if countryCode == "" {
		countryCode = "TH"
	}

	if categoryID == "" {
		return u.songRepo.GetTrends(countryCode)
	}

	catRepo, err := u.resolveCategoryRepo(categoryID)
	if err != nil {
		return nil, err
	}

	videos, err := catRepo.GetVideos(countryCode)
	if err != nil {
		return nil, err
	}

	return mapCategoryVideosToSongs(videos, categoryID), nil
}

func (u *songUsecase) GetNewReleases(countryCode, categoryID string) ([]domain.Song, error) {
	if countryCode == "" {
		countryCode = "TH"
	}

	if categoryID == "" {
		return u.songRepo.GetNewReleases(countryCode)
	}

	catRepo, err := u.resolveCategoryRepo(categoryID)
	if err != nil {
		return nil, err
	}

	videos, err := catRepo.GetNewVideos(countryCode)
	if err != nil {
		return nil, err
	}

	return mapCategoryVideosToSongs(videos, categoryID), nil
}

func (u *songUsecase) GetMVs(countryCode string) ([]domain.Song, error) {
	if countryCode == "" {
		countryCode = "TH"
	}
	return u.songRepo.GetMVs(countryCode)
}
