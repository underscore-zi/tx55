package testclient

import (
	"tx55/pkg/metalgearonline1/handlers/hostgame"
	"tx55/pkg/metalgearonline1/types"
)

func (c *TestClient) HostCreateGame(title, description, password string) (err error) {
	var data types.CreateGameOptions
	copy(data.Name[:], title)
	copy(data.Description[:], description)
	if len(password) > 0 {
		copy(data.Password[:], password)
		data.HasPassword = true
	}
	data.Rules[0] = types.GameRules{Mode: types.ModeTeamDeathmatch, Map: types.MapKillhouseA}
	data.Rules[1] = types.GameRules{Mode: types.ModeTeamDeathmatch, Map: types.MapKillhouseB}
	data.Rules[2] = types.GameRules{Mode: types.ModeTeamDeathmatch, Map: types.MapKillhouseC}

	if err = c.Send(types.ClientCreateGame, data); err != nil {
		return
	}
	if err = c.ExpectErrorCode(0); err != nil {
		return
	}

	if err = c.Send(types.ClientHostReadyToCreate, []byte{}); err != nil {
		return
	}
	// Ready to create gets an empty response
	if _, err = c.Receive(); err != nil {
		return
	}

	return
}

func (c *TestClient) HostPlayerJoin(UserID types.UserID) (err error) {
	args := hostgame.ArgsHostPlayerJoin{UserID: UserID}
	if err = c.Send(types.ClientHostPlayerJoin, args); err != nil {
		return
	}

	if err = c.ExpectErrorCode(0); err != nil {
		return
	}
	return nil
}

func (c *TestClient) HostPlayerLeave(UserID types.UserID) (err error) {
	args := hostgame.ArgsHostPlayerLeave{UserID: UserID}
	if err = c.Send(types.ClientHostPlayerLeave, args); err != nil {
		return
	}

	if err = c.ExpectErrorCode(0); err != nil {
		return
	}
	return nil
}

func (c *TestClient) HostPlayerKicked(UserID types.UserID) (err error) {
	args := hostgame.ArgsHostPlayerKicked{UserID: UserID}
	if err = c.Send(types.ClientHostPlayerKicked, args); err != nil {
		return
	}

	if err = c.ExpectErrorCode(0); err != nil {
		return
	}
	return nil
}

func (c *TestClient) HostPlayerJoinTeam(UserID types.UserID, TeamID types.Team) (err error) {
	args := hostgame.ArgsHostPlayerJoinTeam{UserID: UserID, TeamID: TeamID}
	if err = c.Send(types.ClientHostPlayerJoinTeam, args); err != nil {
		return
	}

	if err = c.ExpectErrorCode(0); err != nil {
		return
	}
	return nil
}

func (c *TestClient) HostPlayerStats(UserID types.UserID, Stats types.HostReportedStats) (err error) {
	args := hostgame.ArgsHostPlayerStats{UserID: UserID, Stats: Stats}
	if err = c.Send(types.ClientHostPlayerStats, args); err != nil {
		return
	}

	if err = c.ExpectErrorCode(0); err != nil {
		return
	}
	return nil
}
