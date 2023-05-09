package konamiserver

import (
	"tx55/pkg/packet"
)

type HookType int
type HookResult int

// HookFunc needs to support two arguments. The first will always be provided and is the packet that triggered the hook
// output hooks will also get a `request` packet that triggered the hook in the first place
type HookFunc func(packet *packet.Packet, request *packet.Packet, out chan packet.Packet) HookResult

const (
	HookBefore HookType = iota
	HookAfter
	HookOutputPacket
)

const (
	HookResultContinue HookResult = iota
	HookResultStop
	// HookResultDrop only implemented for output hooks
	HookResultDrop
)

func (s *Server) AddHook(command uint16, hookType HookType, f HookFunc) {
	switch hookType {
	case HookBefore:
		s.beforeHooks[command] = append(s.beforeHooks[command], f)
	case HookAfter:
		s.afterHooks[command] = append(s.afterHooks[command], f)
	case HookOutputPacket:
		s.outputHooks[command] = append(s.outputHooks[command], f)
	}
}

// ----- Client Functions -----

// processReplies wraps all the replies cause by a single packet, so we can apply output hooks
// it's a bit round-about create a new channel just for that, but it lets us associate the request with the outputs
// hooking later in the writer loop removes that information
func (c *client) processReplies(request packet.Packet, out chan packet.Packet, in chan packet.Packet) {
	for {
		select {
		case msg, ok := <-in:
			if !ok {
				return
			}
			if hooks, found := c.server.outputHooks[msg.Type()]; found {
				switch c.dispatchHooks(hooks, &msg, &request, in) {
				case HookResultDrop:
					continue
				case HookResultStop:
				}
			}
			out <- msg

		}
	}
}

// dispatch will be executed in a new goroutine, so we can immediately process this packet
func (c *client) dispatch(p *packet.Packet) {
	replies := make(chan packet.Packet)
	go c.processReplies(*p, c.out, replies)
	defer close(replies)

	if hooks, found := c.server.beforeHooks[(*p).Type()]; found {
		switch c.dispatchHooks(hooks, p, nil, replies) {
		case HookResultStop:
			return
		}
	}

	if err := c.gameClient.OnPacket(p, replies); err != nil {
		c.once.Do(c.cleanup)
		return
	}

	if hooks, found := c.server.afterHooks[(*p).Type()]; found {
		switch c.dispatchHooks(hooks, p, nil, replies) {
		case HookResultStop:
			return
		}
	}

	return
}

// dispatchHooks is a simple wrapper to iterate the hooks and execute them, stopping when expected
func (c *client) dispatchHooks(hooks []HookFunc, p *packet.Packet, request *packet.Packet, out chan packet.Packet) HookResult {
	for _, hook := range hooks {
		switch hook(p, request, out) {
		case HookResultStop:
			return HookResultStop
		case HookResultDrop:
			return HookResultDrop
		}
	}
	return HookResultContinue
}
