package crons

import (
	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

func Schedule(s *gocron.Scheduler, db *gorm.DB) error {
	if _, err := s.Every(1).Hour().Do(UpdateRankings, db); err != nil {
		return err
	}
	if _, err := s.Every(1).Day().Monday().At("00:00").Do(ClearWeeklyStats, db); err != nil {
		return err
	}

	return nil
}
