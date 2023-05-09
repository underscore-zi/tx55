package auth

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(GetNonceHandler{})
}

type GetNonceHandler struct{}

func (h GetNonceHandler) Type() types.PacketType {
	return types.ClientGetNonce
}

func (h GetNonceHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h GetNonceHandler) Handle(_ *session.Session, _ *packet.Packet) ([]types.Response, error) {
	// Originally they probably changed this every time as it is used for a digest auth mechanism. To do that
	// means storing passwords in a recoverable way both the server can know the correct response without needing to
	// transmit the secret over the wire. To avoid storing the passwords we use a static value here. The trade-off is
	// that we do have to transmit the password over the wire, and we don't have the protection of TLS on game traffic
	return []types.Response{ResponseNonce{
		Nonce: types.NONCE,
	}}, nil
}

// --- Packets ---

type ResponseNonce struct {
	Nonce [16]byte
}

func (r ResponseNonce) Type() types.PacketType { return types.ServerNonce }
