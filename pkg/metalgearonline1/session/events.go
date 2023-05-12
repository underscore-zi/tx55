package session

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"tx55/pkg/metalgearonline1/types"
)

type EventType string

const (
	EventTypeGameCreated      EventType = "game_created"
	EventTypeGameDeleted      EventType = "game_deleted"
	EventTypeGameNewRound     EventType = "game_new_round"
	EventTypeGamePlayerJoined EventType = "game_player_joined"
	EventTypeGamePlayerLeft   EventType = "game_player_left"
)

func (s *Session) publishEvent(e EventType, data interface{}) {
	endpoint, found := os.LookupEnv("EVENTS_ENDPOINT")
	if !found {
		return
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"event": e,
		"lobby": s.LobbyID,
		"data":  data,
	})

	if err != nil {
		s.Log.WithError(err).WithField("event", e).Error("failed to marshal event data")
		return
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.Log.WithError(err).WithField("event", e).Error("failed to post event data")
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		s.Log.WithField("event", e).WithField("status_code", resp.StatusCode).Error("failed to post event data")
	}

	return
}

func (s *Session) EventGameCreated(gameID uint, args *types.CreateGameOptions) {
	type roundEntry struct {
		Map  types.GameMapString
		Mode types.GameModeString
	}

	var rounds []roundEntry
	for _, r := range args.Rules {
		if r.Map == 0 {
			break
		}
		rounds = append(rounds, roundEntry{
			Map:  r.Map.String(),
			Mode: r.Mode.String(),
		})
	}
	go s.publishEvent(EventTypeGameCreated, map[string]interface{}{
		"name":         types.BytesToString(args.Name[:]),
		"has_password": args.HasPassword,
		"game_id":      gameID,
		"host":         types.BytesToString(s.User.DisplayName[:]),
		"rules":        rounds,
	})
}

func (s *Session) EventGameDeleted() {
	go s.publishEvent(EventTypeGameDeleted, map[string]interface{}{
		"game_id": s.GameState.GameID,
	})
}

func (s *Session) EventGameNewRound(round byte) {
	go s.publishEvent(EventTypeGameNewRound, map[string]interface{}{
		"game_id": s.GameState.GameID,
		"round":   round,
		"map":     s.GameState.Rules[round].Map.String(),
		"mode":    s.GameState.Rules[round].Mode.String(),
	})
}

func (s *Session) EventGamePlayerJoined(id types.UserID) {
	go s.publishEvent(EventTypeGamePlayerJoined, map[string]interface{}{
		"game_id": s.GameState.GameID,
		"user_id": id,
	})
}

func (s *Session) EventGamePlayerLeft(id types.UserID) {
	go s.publishEvent(EventTypeGamePlayerLeft, map[string]interface{}{
		"game_id": s.GameState.GameID,
		"user_id": id,
	})
}
