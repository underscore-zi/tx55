package auth

import (
	"reflect"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/models"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

func init() {
	handlers.Register(NewsListHandler{})
}

type NewsListHandler struct{}

func (h NewsListHandler) Type() types.PacketType {
	return types.ClientGetNewsList
}

func (h NewsListHandler) ArgumentTypes() []reflect.Type {
	return []reflect.Type{}
}

func (h NewsListHandler) Handle(sess *session.Session, _ *packet.Packet) ([]types.Response, error) {
	var out []types.Response
	var entries []models.News
	out = append(out, ResponseNewsListStart{})

	_ = sess.DB.Where("topic != ?", "policy").Find(&entries)
	for _, entry := range entries {
		news := ResponseNewsListEntry{
			ID: uint32(entry.ID),
		}
		copy(news.Time[:], entry.Time.Format("2006-01-02 15:04:05"))
		copy(news.Topic[:], entry.Topic)
		copy(news.Body[:], entry.Body)
		out = append(out, news)
	}

	out = append(out, ResponseNewsListEnd{})

	return out, nil
}

type ResponseNewsListStart types.ResponseEmpty

func (r ResponseNewsListStart) Type() types.PacketType { return types.ServerNewsListStart }

type ResponseNewsListEnd types.ResponseEmpty

func (r ResponseNewsListEnd) Type() types.PacketType { return types.ServerNewsListEnd }

type ResponseNewsListEntry struct {
	ID      uint32
	Unknown byte
	Time    [19]byte
	Topic   [64]byte
	Body    [900]byte `packet:"truncate"`
}

func (r ResponseNewsListEntry) Type() types.PacketType { return types.ServerNewsListEntry }
