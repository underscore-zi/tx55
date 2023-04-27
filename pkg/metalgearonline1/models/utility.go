package models

import (
	"errors"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"tx55/pkg/metalgearonline1/types"
)

func addMatchingFields(dest, src interface{}, blacklist []string) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return errors.New("dest must be a non-nil pointer")
	}
	destVal = destVal.Elem()

	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() != reflect.Ptr || srcVal.IsNil() {
		return errors.New("src must be a non-nil pointer")
	}
	srcVal = srcVal.Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Type().Field(i)
		srcFieldName := srcField.Name
		for _, blacklisted := range blacklist {
			if blacklisted == srcFieldName {
				continue
			}
		}

		if strings.HasSuffix(srcFieldName, "ID") {
			continue
		}

		destField := destVal.FieldByName(srcFieldName)
		if !destField.IsValid() || !destField.CanSet() {
			continue
		}

		if destField.Type() == srcField.Type {
			switch srcField.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				destField.SetInt(destField.Int() + srcVal.Field(i).Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				destField.SetUint(destField.Uint() + srcVal.Field(i).Uint())
			case reflect.Float32, reflect.Float64:
				destField.SetFloat(destField.Float() + srcVal.Field(i).Float())
			}
		}
	}
	return nil
}

func GetPlayerStats(db *gorm.DB, UserID types.UserID) (allTimeStats, weeklyStats types.PeriodStats, err error) {
	var stats []PlayerStats

	allTimeStats.Period = types.PeriodAllTime
	weeklyStats.Period = types.PeriodWeekly

	db.Where("user_id = ?", UserID).Find(&stats)
	for _, stat := range stats {
		var target *types.PeriodStats
		var targetMode *types.GameTypeStatsWithRank
		switch stat.Period {
		case types.PeriodArchive:
			fallthrough
		case types.PeriodAllTime:
			target = &allTimeStats
		case types.PeriodWeekly:
			target = &weeklyStats
		default:
			// This should never be reached
			panic("Invalid period")
		}
		switch stat.Mode {
		case types.ModeDeathmatch:
			targetMode = &target.Deathmatch
		case types.ModeTeamDeathmatch:
			targetMode = &target.TeamDeathmatch
		case types.ModeCapture:
			targetMode = &target.Capture
		case types.ModeSneaking:
			targetMode = &target.Sneaking
		case types.ModeRescue:
			targetMode = &target.Rescue
		default:
			// This should never be reached
			panic("Invalid mode")
		}

		if err = addMatchingFields(&targetMode.Stats, &stat, []string{"DeathStreak", "KillStreak"}); err != nil {
			return
		}

		if stat.KillStreak > targetMode.Stats.KillStreak {
			targetMode.Stats.KillStreak = stat.KillStreak
		}
		if stat.DeathStreak > targetMode.Stats.DeathStreak {
			targetMode.Stats.DeathStreak = stat.DeathStreak
		}

	}

	return
}
