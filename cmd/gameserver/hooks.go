package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"time"
	"tx55/pkg/konamiserver"
	"tx55/pkg/metalgearonline1/handlers/auth"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

// ALl the code here is really just a hack to be able to transfer sessions from the old (live) game server to the new database
// so people can connect to this new one without having to get rid of the original server

var OriginalDb *sql.DB
var GormDb *gorm.DB

func connectToDatabase() *sql.DB {
	oldDSN, _ := os.LookupEnv("OLD_DSN")
	if oldDSN == "" {
		return nil
	}

	db, err := sql.Open("mysql", oldDSN)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(2 * time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	return db
}

func init() {
	OriginalDb = connectToDatabase()
}

func b64ToGameOptions(source string) (options types.CreateGameOptions, err error) {
	decoded, err := base64.StdEncoding.DecodeString(source)
	if err != nil {
		return
	}
	err = binary.Read(bytes.NewReader(decoded), binary.BigEndian, &options)
	return
}

func b64ToPlayerSettings(source string) (settings types.PlayerSpecificSettings, err error) {
	decoded, err := base64.StdEncoding.DecodeString(source)
	if err != nil {
		return
	}
	err = binary.Read(bytes.NewReader(decoded), binary.BigEndian, &settings)
	return
}

func hookLogin(p, _ *packet.Packet, _ chan packet.Packet) konamiserver.HookResult {
	args := auth.ArgsLoginSession{}
	if int((*p).Length()) != binary.Size(args) {
		return konamiserver.HookResultContinue
	}
	if err := (*p).DataInto(&args); err != nil {
		l.WithError(err).Error("Failed unmarshalling login (session) packet")
		return konamiserver.HookResultContinue
	}

	uid, err := uuid.FromBytes(args.SessionID[:])
	if err != nil {
		l.WithError(err).Error("Failed parsing session id")
		return konamiserver.HookResultContinue
	}
	s := &models.Session{ID: uid}
	if tx := GormDb.Model(s).First(s); tx.Error == nil {
		// The session exists, we don't need to do anything with this one
		return konamiserver.HookResultContinue
	} else if tx.Error != gorm.ErrRecordNotFound {
		l.WithError(tx.Error).Error("Failed looking up session")
		return konamiserver.HookResultContinue
	}

	// Record wasn't found, so we must need to transfer this session from the old database
	// First, lets check if the session exists in the old database
	sid := fmt.Sprintf("%x", args.SessionID)
	query := "SELECT username, displayname, flags, emblem_text, user_settings, game_settings FROM users WHERE session_id = ?"
	rows, err := OriginalDb.Query(query, sid)
	if err != nil {
		l.WithError(err).Error("Couldn't find session in old db")
		return konamiserver.HookResultContinue
	}
	defer func() { _ = rows.Close() }()
	if !rows.Next() {
		l.Error("Unexpected failure of Next()")
		return konamiserver.HookResultContinue
	}

	var flags int
	var username, displayName, userSettings, gameSettings, emblemText string
	if err := rows.Scan(&username, &displayName, &flags, &emblemText, &userSettings, &gameSettings); err != nil {
		l.WithError(err).Error("Failed scanning session")
		return konamiserver.HookResultContinue
	}

	// Now we have the base user information, do we need to create a transfer data?
	gs := &models.GameOptions{}
	u := &models.User{}
	tx := GormDb.Where("username LIKE ?", username).First(u)
	saveGameSettings := false

	if tx.Error == gorm.ErrRecordNotFound {
		// New user time!
		u.Password = "ABC132"
		u.Username = []byte(username)
		u.DisplayName = []byte(displayName)
		if flags&0x10000 > 0 {
			u.HasEmblem = true
			u.EmblemText = []byte(emblemText)
		}
		ps := &models.PlayerSettings{
			UserID: u.ID,
		}

		// Player settings
		rawSettings, err := b64ToPlayerSettings(userSettings)
		if err != nil {
			l.WithError(err).Error("Failed parsing player settings")
			return konamiserver.HookResultContinue
		}
		ps.FromBitfield(rawSettings.Settings)
		ps.FKey0 = rawSettings.FKeys[0][:]
		ps.FKey1 = rawSettings.FKeys[1][:]
		ps.FKey2 = rawSettings.FKeys[2][:]
		ps.FKey3 = rawSettings.FKeys[3][:]
		ps.FKey4 = rawSettings.FKeys[4][:]
		ps.FKey5 = rawSettings.FKeys[5][:]
		ps.FKey6 = rawSettings.FKeys[6][:]
		ps.FKey7 = rawSettings.FKeys[7][:]
		ps.FKey8 = rawSettings.FKeys[8][:]
		ps.FKey9 = rawSettings.FKeys[9][:]
		ps.FKey10 = rawSettings.FKeys[10][:]
		ps.FKey11 = rawSettings.FKeys[11][:]
		u.PlayerSettings = *ps

		// Create Game Settings
		if gameSettings != "" {
			gameOpts, err := b64ToGameOptions(gameSettings)
			if err != nil {
				l.WithError(err).Error("Failed parsing game settings")
				return konamiserver.HookResultContinue
			}
			gs.FromCreateGameOptions(&gameOpts)
			saveGameSettings = true
		}
	}

	_ = GormDb.Transaction(func(tx *gorm.DB) error {

		if tx = GormDb.Save(u); tx.Error != nil {
			l.WithError(tx.Error).Error("Failed saving user")
			return tx.Error
		}

		if saveGameSettings {
			gs.UserID = u.ID
			if tx = GormDb.Save(gs); tx.Error != nil {
				l.WithError(tx.Error).Error("Failed saving game settings")
			}
		}

		s.UserID = u.ID
		if tx = GormDb.Save(s); tx.Error != nil {
			l.WithError(tx.Error).Error("Failed saving session")
		}
		return tx.Error
	})

	return konamiserver.HookResultContinue
}

func hookConnectionInfo(p, _ *packet.Packet, _ chan packet.Packet) konamiserver.HookResult {
	if ip, found := os.LookupEnv("FORCED_HOST_REMOTE_ADDR"); found {
		if len(ip) >= 15 {
			ip = ip[:15]
		} else {
			ip = ip + "\x00"
		}
		data := (*p).Data()
		copy(data[4:], ip)
		(*p).SetData(data)
	}
	return konamiserver.HookResultContinue
}

//goland:noinspection GoUnusedFunction
func hookResponseFromFile(s *konamiserver.Server, cmd types.PacketType, outCmd []types.PacketType, outFile []string) {
	if len(outCmd) != len(outFile) {
		panic("outCmd and outFile must be the same length")
	}

	s.AddHook(uint16(cmd), konamiserver.HookBefore, func(p, req *packet.Packet, out chan packet.Packet) konamiserver.HookResult {
		for i, file := range outFile {
			var bs []byte
			var err error

			if file != "" {
				bs, err = os.ReadFile(file)
				if err != nil {
					log.WithError(err).WithField("file", file).Error("Failed reading file")
					break
				}
			}

			p := packet.New()
			p.SetType(uint16(outCmd[i]))
			p.SetData(bs)
			out <- p
		}
		return konamiserver.HookResultStop
	})
}
