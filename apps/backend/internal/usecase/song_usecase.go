package usecase

import (
	"log"

	"tunetrend-backend/internal/domain"
)

type songUsecase struct {
	songRepo domain.SongRepository
	ytRepo   domain.YouTubeRepository
}

func NewSongUsecase(songRepo domain.SongRepository, ytRepo domain.YouTubeRepository) domain.SongUsecase {
	return &songUsecase{
		songRepo: songRepo,
		ytRepo:   ytRepo,
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

func (u *songUsecase) GetTrends(countryCode string) ([]domain.Song, error) {
	if countryCode == "" {
		countryCode = "TH"
	}

	return u.songRepo.GetTrends(countryCode)
}

func (u *songUsecase) GetNewReleases(countryCode string) ([]domain.Song, error) {
	if countryCode == "" {
		countryCode = "TH"
	}
	return u.songRepo.GetNewReleases(countryCode)
}

func (u *songUsecase) GetMVs(countryCode string) ([]domain.Song, error) {
	if countryCode == "" {
		countryCode = "TH"
	}
	return u.songRepo.GetMVs(countryCode)
}
