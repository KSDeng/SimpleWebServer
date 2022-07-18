package packet

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	CommandConn = iota + 0x01
	CommandLogin
	CommandUpdatePhotoPath
	CommandGetUserData
	CommandUpdateNickname
)

const (
	CommandConnAck = iota + 0x80
	CommandLoginAck
	CommandUpdatePhotoPathAck
	CommandGetUserDataAck
	CommandUpdateNicknameAck
)

type Packet interface {
	Decode([]byte) error     // []byte -> struct
	Encode() ([]byte, error) // struct -> []byte
}

/* Login request begins */
type LoginRequestPacket struct {
	ID         string
	PacketBody LoginRequestPacketBody
}

func (p *LoginRequestPacket) Decode(pktBody []byte) error {
	p.ID = string(pktBody[:8])
	err := json.Unmarshal(pktBody[8:], &p.PacketBody)
	if err != nil {
		return err
	}
	return nil
}

func (p *LoginRequestPacket) Encode() ([]byte, error) {
	packetBodyBytes, _ := json.Marshal(p.PacketBody)
	return bytes.Join([][]byte{[]byte(p.ID[:8]), packetBodyBytes}, nil), nil
}

/* Login request ends */

/* Login reply begins */
type LoginRequestAck struct {
	ID         string
	PacketBody LoginRequestAckBody
}

func (p *LoginRequestAck) Decode(pktBody []byte) error {
	p.ID = string(pktBody[:8])
	err := json.Unmarshal(pktBody[8:], &p.PacketBody)
	if err != nil {
		return err
	}
	return nil
}

func (p *LoginRequestAck) Encode() ([]byte, error) {
	packetBodyBytes, _ := json.Marshal(p.PacketBody)
	return bytes.Join([][]byte{[]byte(p.ID[:8]), packetBodyBytes}, nil), nil
}

/* Login reply ends */

/* Update photo path request begins */
type UpdatePhotoPathPacket struct {
	ID         string
	PacketBody UpdatePhotoPathPacketBody
}

func (p *UpdatePhotoPathPacket) Decode(pktBody []byte) error {
	p.ID = string(pktBody[:8])
	err := json.Unmarshal(pktBody[8:], &p.PacketBody)
	if err != nil {
		return err
	}
	return nil
}

func (p *UpdatePhotoPathPacket) Encode() ([]byte, error) {
	packetBodyBytes, _ := json.Marshal(p.PacketBody)
	return bytes.Join([][]byte{[]byte(p.ID[:8]), packetBodyBytes}, nil), nil
}

/* Update photo path request ends */

/* Update photo path reply begins */
type UpdatePhotoPathAck struct {
	ID         string
	PacketBody UpdatePhotoPathAckBody
}

func (p *UpdatePhotoPathAck) Decode(pktBody []byte) error {
	p.ID = string(pktBody[:8])
	err := json.Unmarshal(pktBody[8:], &p.PacketBody)
	if err != nil {
		return err
	}
	return nil
}

func (p *UpdatePhotoPathAck) Encode() ([]byte, error) {
	packetBodyBytes, _ := json.Marshal(p.PacketBody)
	return bytes.Join([][]byte{[]byte(p.ID[:8]), packetBodyBytes}, nil), nil
}

/* Update photo path reply ends */

/* Update nickname begins */
type UpdateNicknameRequest struct {
	ID         string
	PacketBody UpdateNicknamePacketBody
}

func (p *UpdateNicknameRequest) Decode(pktBody []byte) error {
	p.ID = string(pktBody[:8])
	err := json.Unmarshal(pktBody[8:], &p.PacketBody)
	if err != nil {
		return err
	}
	return nil
}

func (p *UpdateNicknameRequest) Encode() ([]byte, error) {
	packetBodyBytes, _ := json.Marshal(p.PacketBody)
	return bytes.Join([][]byte{[]byte(p.ID[:8]), packetBodyBytes}, nil), nil
}

/* Update nickname ends */

/* Update nickname ack begins */
type UpdateNicknameAck struct {
	ID         string
	PacketBody UpdateNicknameAckBody
}

func (p *UpdateNicknameAck) Decode(pktBody []byte) error {
	p.ID = string(pktBody[:8])
	err := json.Unmarshal(pktBody[8:], &p.PacketBody)
	if err != nil {
		return err
	}
	return nil
}

func (p *UpdateNicknameAck) Encode() ([]byte, error) {
	packetBodyBytes, _ := json.Marshal(p.PacketBody)
	return bytes.Join([][]byte{[]byte(p.ID[:8]), packetBodyBytes}, nil), nil
}

/* Update nickname ack ends */

/* Get user data request begins */
type GetUserDataRequestPacket struct {
	ID         string
	PacketBody GetUserDataRequestBody
}

func (p *GetUserDataRequestPacket) Decode(pktBody []byte) error {
	p.ID = string(pktBody[:8])
	err := json.Unmarshal(pktBody[8:], &p.PacketBody)
	if err != nil {
		return err
	}
	return nil
}

func (p *GetUserDataRequestPacket) Encode() ([]byte, error) {
	packetBodyBytes, _ := json.Marshal(p.PacketBody)
	return bytes.Join([][]byte{[]byte(p.ID[:8]), packetBodyBytes}, nil), nil
}

/* Get user data request ends */

/* Get user data reply begins */
type GetUserDataAck struct {
	ID         string
	PacketBody GetUserDataAckBody
}

func (p *GetUserDataAck) Decode(pktBody []byte) error {
	p.ID = string(pktBody[:8])
	err := json.Unmarshal(pktBody[8:], &p.PacketBody)
	if err != nil {
		return err
	}
	return nil
}

func (p *GetUserDataAck) Encode() ([]byte, error) {
	packetBodyBytes, _ := json.Marshal(p.PacketBody)
	return bytes.Join([][]byte{[]byte(p.ID[:8]), packetBodyBytes}, nil), nil
}

/* Get user data reply ends */

func Decode(packet []byte) (Packet, error) {
	commandID := packet[0]
	pktBody := packet[1:]

	switch commandID {
	case CommandConn:
		return nil, nil
	case CommandConnAck:
		return nil, nil
	case CommandLogin:
		p := LoginRequestPacket{}
		err := p.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &p, nil
	case CommandLoginAck:
		p := LoginRequestAck{}
		err := p.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &p, nil
	case CommandUpdatePhotoPath:
		p := UpdatePhotoPathPacket{}
		err := p.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &p, nil
	case CommandUpdatePhotoPathAck:
		p := UpdatePhotoPathAck{}
		err := p.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &p, nil
	case CommandGetUserData:
		p := GetUserDataRequestPacket{}
		err := p.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &p, nil
	case CommandGetUserDataAck:
		p := GetUserDataAck{}
		err := p.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &p, nil
	case CommandUpdateNickname:
		p := UpdateNicknameRequest{}
		err := p.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &p, nil
	case CommandUpdateNicknameAck:
		p := UpdateNicknameAck{}
		err := p.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &p, nil
	default:
		return nil, fmt.Errorf("unknown commandID [%d]", commandID)
	}
}

func Encode(p Packet) ([]byte, error) {
	var commandID uint8
	var pktBody []byte
	var err error

	switch t := p.(type) {
	case *LoginRequestPacket:
		commandID = CommandLogin
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *LoginRequestAck:
		commandID = CommandLoginAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *UpdatePhotoPathPacket:
		commandID = CommandUpdatePhotoPath
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *UpdatePhotoPathAck:
		commandID = CommandUpdatePhotoPathAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *GetUserDataRequestPacket:
		commandID = CommandGetUserData
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *GetUserDataAck:
		commandID = CommandGetUserDataAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *UpdateNicknameRequest:
		commandID = CommandUpdateNickname
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *UpdateNicknameAck:
		commandID = CommandUpdateNicknameAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type [%s]", t)
	}
	return bytes.Join([][]byte{[]byte{commandID}, pktBody}, nil), nil
}
