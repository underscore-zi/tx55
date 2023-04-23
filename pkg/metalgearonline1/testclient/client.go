package testclient

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net"
	"tx55/pkg/metalgearonline1/handlers/auth"
	"tx55/pkg/metalgearonline1/handlers/lobby"
	"tx55/pkg/metalgearonline1/types"
)

type TestClient struct {
	Key  []byte
	seq  uint32
	conn net.Conn
}

func (c *TestClient) LoginWithSession(sessionID string) error {
	uid, _ := uuid.Parse(sessionID)
	data := auth.ArgsLoginSession{}
	data.SessionID = [16]byte(uid[:])

	_ = c.Send(types.ClientLogin, data)

	resp, err := c.Receive()
	if err != nil {
		return err
	}

	responseData := (*resp).Data()
	if len(responseData) < 4 {
		return errors.New("response too small")
	}

	// convert responseData to uint32
	var errorCode uint32
	errorCode |= uint32(responseData[0]) << 24
	errorCode |= uint32(responseData[1]) << 16
	errorCode |= uint32(responseData[2]) << 8
	errorCode |= uint32(responseData[3])

	if errorCode != 0 {
		return errors.New(fmt.Sprintf("code(%d)", int32(errorCode)))
	}

	var payload auth.ResponseLogin
	err = (*resp).DataInto(&payload)
	if err != nil {
		return err
	}

	responseID, _ := uuid.FromBytes(payload.SessionID[:])
	if responseID.String() != sessionID {
		return errors.New("unexpected session ID")
	}

	return nil
}

func (c *TestClient) ReportConnectionInfo(localAddr string, remotePort, localPort uint16) (err error) {
	data := lobby.ArgsReportConnectionInfo{
		RemotePort: remotePort,
		LocalPort:  localPort,
	}
	copy(data.LocalAddr[:], localAddr)
	if err = c.Send(types.ClientReportConnectionInfo, data); err != nil {
		return
	}

	if err = c.ExpectErrorCode(0); err != nil {
		return
	}

	return
}
