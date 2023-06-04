package handlers

import (
	"encoding/binary"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strings"
	"tx55/pkg/metalgearonline1/session"
	"tx55/pkg/metalgearonline1/types"
	"tx55/pkg/packet"
)

var AllHandlers = map[types.PacketType]Handler{}

// Handler is the interface that all handlers must implement
type Handler interface {
	// Type is a singular packet type that the handler will handle
	Type() types.PacketType
	// ArgumentTypes is a list of argument types the handler expects. Argument types must be fixed size structs
	ArgumentTypes() []reflect.Type
	// Handle is the default handler that'll be called with teh session and packet
	// However you can have other `Handle*` methods. Instead of a packet these can take a pointer to one of the types
	// returned from ArgumentTypes(). The dispatcher can automatically discover these methods and call them
	Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error)
}

// Register should be called from the init() function of each handler package
func Register(h Handler) {
	if existing, ok := AllHandlers[h.Type()]; ok {
		existingName := reflect.TypeOf(existing).String()
		typeName := reflect.TypeOf(h).String()
		fmt.Printf("[0x%04x] Could not register %s (0x%04x: %s)\n", h.Type(), typeName, typeName, existingName)
		os.Exit(1)
	}

	AllHandlers[h.Type()] = h
}

// getArgs compares the data length with the expected argument struct sizes and returns the first match
// for MGO1 it's the case that when the args are different, the size is different
func getArgs(h Handler, p *packet.Packet) (reflect.Type, any, error) {
	dataLen := (*p).Length()
	argOptions := h.ArgumentTypes()
	for _, T := range argOptions {
		impl := reflect.New(T).Interface()
		implSize := uint16(binary.Size(impl))
		sizesMatch := dataLen == implSize

		if sizesMatch {
			if err := (*p).DataInto(impl); err != nil {
				return nil, nil, err
			}
			return T, impl, nil
		}
	}

	return nil, nil, nil
}

// Handle will dispatch automatically find the correct handler for a given packet
// If it can find a handler with arguments that match what was received it'll send
// the packet to that handler, otherwise it defaults to the interface Handler() func
func Handle(sess *session.Session, packet *packet.Packet) ([]types.Response, error) {
	cmd := types.PacketType((*packet).Type())
	handler, ok := AllHandlers[cmd]
	if !ok {
		return nil, ErrHandlerNotFound
	}

	argT, argV, err := getArgs(handler, packet)
	if err != nil {
		return nil, err
	}
	if argT != nil {
		handlerType := reflect.TypeOf(handler)
		for i := 0; i < handlerType.NumMethod(); i++ {
			method := handlerType.Method(i)
			if !strings.HasPrefix(method.Name, "Handle") {
				continue
			}
			if method.Type.NumIn() != 3 {
				continue
			}
			thirdArg := method.Type.In(2)
			if thirdArg != reflect.PointerTo(argT) {
				continue
			}
			args := []reflect.Value{
				reflect.ValueOf(handler),
				reflect.ValueOf(sess),
				reflect.ValueOf(argV),
			}
			returns := method.Func.Call(args)

			var replies []types.Response
			var err error

			// Apparently nils can Interface() but then fail to cast
			if returns[0].Interface() != nil {
				replies = returns[0].Interface().([]types.Response)
			}

			if returns[1].Interface() != nil {
				err = returns[1].Interface().(error)
			}

			return replies, err
		}
	}

	if len(handler.ArgumentTypes()) > 0 {
		sess.LogEntry().WithFields(log.Fields{
			"type":    (*packet).Type(),
			"handler": reflect.TypeOf(handler).String(),
			"len":     (*packet).Length(),
		}).Error("No argument-specific handler found")
	}

	return handler.Handle(sess, packet)
}
