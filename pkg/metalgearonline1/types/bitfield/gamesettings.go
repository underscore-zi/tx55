package bitfield

type GameSettings struct {
	Bytes [2]Bitfield
}

const (
	// First Byte
	UniqueCharactersPos = 0
	RumbleRosesGirlsPos = 1
	GhostsPos           = 3
	FriendlyFirePos     = 4
	TeamAutoBalancePos  = 6
	IdleKickPos         = 7
	// Second Byte
	TeamKillKickPos = 0
	VoiceChatPos    = 1
)

func (b *GameSettings) GetUniqueCharacters() bool {
	return b.Bytes[0].GetBit(UniqueCharactersPos)
}

func (b *GameSettings) SetUniqueCharacters(val bool) {
	b.Bytes[0].SetBit(UniqueCharactersPos, val)
}

func (b *GameSettings) GetRumbleRosesGirls() bool {
	return b.Bytes[0].GetBit(RumbleRosesGirlsPos)
}

func (b *GameSettings) SetRumbleRosesGirls(val bool) {
	b.Bytes[0].SetBit(RumbleRosesGirlsPos, val)
}

func (b *GameSettings) GetGhosts() bool {
	return b.Bytes[0].GetBit(GhostsPos)
}

func (b *GameSettings) SetGhosts(val bool) {
	b.Bytes[0].SetBit(GhostsPos, val)
}

func (b *GameSettings) GetFriendlyFire() bool {
	return b.Bytes[0].GetBit(FriendlyFirePos)
}

func (b *GameSettings) SetFriendlyFire(val bool) {
	b.Bytes[0].SetBit(FriendlyFirePos, val)
}

func (b *GameSettings) GetTeamAutoBalance() bool {
	return b.Bytes[0].GetBit(TeamAutoBalancePos)
}

func (b *GameSettings) SetTeamAutoBalance(val bool) {
	b.Bytes[0].SetBit(TeamAutoBalancePos, val)
}

func (b *GameSettings) GetIdleKick() bool {
	return b.Bytes[0].GetBit(IdleKickPos)
}

func (b *GameSettings) SetIdleKick(val bool) {
	b.Bytes[0].SetBit(IdleKickPos, val)
}

func (b *GameSettings) GetTeamKillKick() bool {
	return b.Bytes[1].GetBit(TeamKillKickPos)
}

func (b *GameSettings) SetTeamKillKick(val bool) {
	b.Bytes[1].SetBit(TeamKillKickPos, val)
}

func (b *GameSettings) GetVoiceChat() bool {
	return b.Bytes[1].GetBit(VoiceChatPos)
}

func (b *GameSettings) SetVoiceChat(val bool) {
	b.Bytes[1].SetBit(VoiceChatPos, val)
}
