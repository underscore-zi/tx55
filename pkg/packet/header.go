package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func (h RawHeader) Bytes() []byte {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, h); err != nil {
		log.WithFields(log.Fields{
			"type": fmt.Sprintf("%04x", h.Cmd),
			"seq":  h.Seq,
			"len":  h.Len,
			"md5":  fmt.Sprintf("%x", h.Md5),
		}).WithError(err).Error("Unable to write header to buffer")
	}
	return buf.Bytes()
}
