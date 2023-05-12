package restapi

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
)

func ToLobbyJSON(lobby models.Lobby) LobbyJSON {
	return LobbyJSON{
		ID:      lobby.ID,
		Name:    lobby.Name,
		Players: lobby.Players,
	}
}

func ToUserJSON(user *models.User) *UserJSON {
	if user == nil || user.ID == 0 {
		return nil
	}
	return &UserJSON{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt.Time,
		DisplayName: types.BytesToString(user.DisplayName),
		HasEmblem:   user.HasEmblem,
		EmblemText:  types.BytesToString(user.EmblemText),
		OverallRank: user.OverallRank,
		WeeklyRank:  user.WeeklyRank,
	}
}

func ToPlayerStatsJSON(stats models.PlayerStats) PlayerStatsJSON {
	return PlayerStatsJSON{
		UserID:       stats.UserID,
		UpdatedAt:    stats.UpdatedAt,
		Period:       stats.Period,
		PeriodString: stats.Period.String(),
		Mode:         stats.Mode,
		ModeString:   stats.Mode.String(),
		Map:          stats.Map,
		MapString:    stats.Map.String(),

		Rank: stats.Rank,

		Kills:              stats.Kills,
		Deaths:             stats.Deaths,
		KillStreak:         stats.KillStreak,
		DeathStreak:        stats.DeathStreak,
		Suicides:           stats.Suicides,
		SelfStuns:          stats.SelfStuns,
		Stuns:              stats.Stuns,
		StunsReceived:      stats.StunsReceived,
		SnakeFrags:         stats.SnakeFrags,
		Points:             stats.Points,
		TeamKills:          stats.TeamKills,
		TeamStuns:          stats.TeamStuns,
		RoundsPlayed:       stats.RoundsPlayed,
		RoundsNoDeath:      stats.RoundsNoDeath,
		KerotansForWin:     stats.KerotansForWin,
		KerotansPlaced:     stats.KerotansPlaced,
		RadioUses:          stats.RadioUses,
		TextChatUses:       stats.TextChatUses,
		CQCAttacks:         stats.CQCAttacks,
		CQCAttacksReceived: stats.CQCAttacksReceived,
		HeadShots:          stats.HeadShots,
		HeadShotsReceived:  stats.HeadShotsReceived,
		TeamWins:           stats.TeamWins,
		KillsWithScorpion:  stats.KillsWithScorpion,
		KillsWithKnife:     stats.KillsWithKnife,
		TimesEaten:         stats.TimesEaten,
		Rolls:              stats.Rolls,
		InfraredGoggleUses: time.Duration(stats.InfraredGoggleUses) * time.Second,
		PlayTime:           time.Duration(stats.PlayTime) * time.Second,
	}
}

func ToGameJSON(game models.Game) GameJSON {
	players := make([]GamePlayersJSON, len(game.Players))
	for i, player := range game.Players {
		players[i] = GamePlayersJSON{
			CreatedAt: player.CreatedAt,
			UpdatedAt: player.UpdatedAt,
			DeletedAt: player.DeletedAt.Time,
			UserID:    player.UserID,
			User:      ToUserJSON(&player.User),
			Team:      player.Team,
			Kills:     player.Kills,
			Deaths:    player.Deaths,
			Score:     player.Score,
			Ping:      player.Ping,
			WasKicked: player.WasKicked,
		}

		switch game.GameOptions.Rules[game.CurrentRound].Mode {
		// TODO: Check if Sneaking does anything special with team joins
		case types.ModeDeathmatch:
			players[i].TeamString = player.Team.UniformString()
		default:
			players[i].TeamString = player.Team.ColorString()
		}
	}

	return GameJSON{
		ID:           game.ID,
		CreatedAt:    game.CreatedAt,
		UpdatedAt:    game.UpdatedAt,
		DeletedAt:    game.DeletedAt.Time,
		LobbyID:      game.LobbyID,
		UserID:       game.UserID,
		Options:      ToGameOptionsJSON(game.GameOptions),
		Players:      players,
		CurrentRound: game.CurrentRound,
	}
}

func ToGameOptionsJSON(opts models.GameOptions) GameOptionsJSON {
	out := GameOptionsJSON{
		Name:              types.BytesToString(opts.Name),
		Description:       types.BytesToString(opts.Description),
		HasPassword:       opts.HasPassword,
		IsHostOnly:        opts.IsHostOnly,
		RedTeam:           opts.RedTeam,
		BlueTeam:          opts.BlueTeam,
		WeaponRestriction: opts.WeaponRestriction,
		MaxPlayers:        opts.MaxPlayers,
		RatingRestriction: opts.RatingRestriction,
		Rating:            opts.Rating,
		SneMinutes:        opts.SneMinutes,
		SneRounds:         opts.SneRounds,
		CapMinutes:        opts.CapMinutes,
		CapRounds:         opts.CapRounds,
		ResMinutes:        opts.ResMinutes,
		ResRounds:         opts.ResRounds,
		TDMMinutes:        opts.TDMMinutes,
		TDMRounds:         opts.TDMRounds,
		TDMTickets:        opts.TDMTickets,
		DMMinutes:         opts.DMMinutes,
		IdleKick:          opts.Bitfield.GetIdleKick(),
		IdleKickMinutes:   opts.IdleKickMinutes,
		TeamKillKick:      opts.Bitfield.GetTeamKillKick(),
		TeamKillCount:     opts.TeamKillCount,
		AutoBalanced:      opts.Bitfield.GetTeamAutoBalance(),
		AutoBalanceCount:  opts.AutoBalance,
		UniqueCharacters:  opts.Bitfield.GetUniqueCharacters(),
		RumbleRoses:       opts.Bitfield.GetRumbleRosesGirls(),
		Ghosts:            opts.Bitfield.GetGhosts(),
		FriendFire:        opts.Bitfield.GetFriendlyFire(),
		HasVoiceChat:      opts.Bitfield.GetVoiceChat(),
	}

	// The rules list is always 15 elements, so we need to truncate it for the JSON
	// This is done by looking for the "null" rules, there is no map 0
	ruleCount := 0
	for _, rule := range opts.Rules {
		if rule.Map == 0 {
			break
		}
		ruleCount++
	}

	out.Rules = make([]GameRuleJSON, ruleCount)
	for i, rule := range opts.Rules[:ruleCount] {
		out.Rules[i] = GameRuleJSON{
			Map:        rule.Map,
			Mode:       rule.Mode,
			MapString:  rule.Map.String(),
			ModeString: rule.Mode.String(),
		}
	}
	return out
}

func ToUserSettingsJSON(settings models.PlayerSettings) UserSettingsJSON {
	return UserSettingsJSON{
		ShowNameTags:         settings.ShowNameTags,
		SwitchSpeed:          uint(settings.SwitchSpeed + 1),
		FPVVertical:          settings.FPVVertical,
		FPVHorizontal:        settings.FPVHorizontal,
		FPVSwitchOrientation: types.SwitchOrientation(settings.FPVSwitchOrientation).String(),
		TPVVertical:          settings.TPVVertical,
		TPVHorizontal:        settings.TPVHorizontal,
		TPVChase:             settings.TPVChase,
		FPVRotationSpeed:     uint(settings.FPVRotationSpeed + 1),
		EquipmentSwitchStyle: types.GearSwitchMode(settings.EquipmentSwitchStyle).String(),
		TPVRotationSpeed:     uint(settings.TPVRotationSpeed + 1),
		WeaponSwitchStyle:    types.GearSwitchMode(settings.WeaponSwitchStyle).String(),
		FKeys: [12]string{
			types.BytesToString(settings.FKey0),
			types.BytesToString(settings.FKey1),
			types.BytesToString(settings.FKey2),
			types.BytesToString(settings.FKey3),
			types.BytesToString(settings.FKey4),
			types.BytesToString(settings.FKey5),
			types.BytesToString(settings.FKey6),
			types.BytesToString(settings.FKey7),
			types.BytesToString(settings.FKey8),
			types.BytesToString(settings.FKey9),
			types.BytesToString(settings.FKey10),
			types.BytesToString(settings.FKey11),
		},
	}
}

func ParamAsInt(c *gin.Context, name string, value int) int {
	param := c.Param(name)
	if param == "" {
		return value
	}
	paramval, err := strconv.Atoi(param)
	if err != nil {
		return value
	}
	return paramval
}
