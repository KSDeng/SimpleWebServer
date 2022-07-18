package main

import (
	"entry-task/config"
	"entry-task/network_potocols/frame"
	"entry-task/network_potocols/packet"
	"entry-task/tcpserver/db"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

var configs config.Configurations

/* Utils begin */
func loadConfigs() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../config")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configs)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}

func checkAndStartMySQLServer() {
	cmd := exec.Command("mysql.server", "status")
	stdout, _ := cmd.Output()
	output_str := string(stdout)

	if strings.Contains(output_str, "not running") {
		cmd_start := exec.Command("mysql.server", "start")
		cmd_start.Run()
	}
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/* Utils end */

/* Handlers begin */
func handleSocketConnection(conn net.Conn) {
	defer conn.Close()

	frameCodec := frame.NewMyFrameCodec()

	for {
		// read from the connection
		// decode the frame to get the payload
		// the payload is undecoded packet

		framePayload, err := frameCodec.Decode(conn)
		if err != nil {
			return
		}

		ackFramePayload, err := handlePacket(framePayload)
		if err != nil {
			fmt.Println("handleConn: handle packet error:", err)
			return
		}

		// write ack frame to the connection
		err = frameCodec.Encode(conn, ackFramePayload)
		if err != nil {
			fmt.Println("handleConn: frame encode error:", err)
			return
		}
	}
}

func handlePacket(framePayload []byte) (ackFramePayload []byte, err error) {
	var p packet.Packet
	p, err = packet.Decode(framePayload)

	if err != nil {
		fmt.Println("handlePacket: packet decode error:", err)
		return
	}

	switch p.(type) {
	case *packet.LoginRequestPacket:
		{
			loginRequest := p.(*packet.LoginRequestPacket)
			fmt.Println("recv login request: ", loginRequest)

			nickname := loginRequest.PacketBody.Nickname

			user_info, err := db.UserInfoByNickname(nickname)
			if err != nil {
				fmt.Println("username or password incorrect")
				loginAckBody := packet.LoginRequestAckBody{
					nickname,
					"",
					"",
					"",
					false,
				}
				loginAck := &packet.LoginRequestAck{
					ID:         loginRequest.ID,
					PacketBody: loginAckBody,
				}
				ackFramePayload, err = packet.Encode(loginAck)
				if err != nil {
					fmt.Println("handlePacket: packet encode error:", err)
					return nil, err
				}
				return ackFramePayload, nil
			}

			password_DB := user_info.Password
			photoPath_DB := user_info.PhotoPath
			sessionId := user_info.Id

			if CheckPasswordHash(loginRequest.PacketBody.Password, password_DB) {
				fmt.Println("password correct")

				loginAckBody := packet.LoginRequestAckBody{
					nickname,
					password_DB,
					photoPath_DB,
					strconv.Itoa(int(sessionId)),
					true,
				}
				loginAck := &packet.LoginRequestAck{
					ID:         loginRequest.ID,
					PacketBody: loginAckBody,
				}
				ackFramePayload, err = packet.Encode(loginAck)
				if err != nil {
					fmt.Println("handlePacket: packet encode error:", err)
					return nil, err
				}

			} else {
				// username or password incorrect
				fmt.Println("username or password incorrect")
				loginAckBody := packet.LoginRequestAckBody{
					nickname,
					"",
					"",
					"",
					false,
				}
				loginAck := &packet.LoginRequestAck{
					ID:         loginRequest.ID,
					PacketBody: loginAckBody,
				}
				ackFramePayload, err = packet.Encode(loginAck)
				if err != nil {
					fmt.Println("handlePacket: packet encode error:", err)
					return nil, err
				}
			}
			return ackFramePayload, nil
		}
	case *packet.UpdatePhotoPathPacket:
		{
			fmt.Println("updating photo path...")
			updateRequest := p.(*packet.UpdatePhotoPathPacket)
			fmt.Println("recv update path request", updateRequest)

			nickname := updateRequest.PacketBody.Nickname
			sessionId := updateRequest.PacketBody.SessionId

			userId, err := strconv.Atoi(sessionId)
			if err != nil {
				fmt.Println(err)
			}
			err = db.UpdatePhotoPathById(int64(userId), updateRequest.PacketBody.PhotoPath)

			updateAckBody := packet.UpdatePhotoPathAckBody{
				nickname,
				err == nil,
			}
			updateAck := &packet.UpdatePhotoPathAck{
				ID:         updateRequest.ID,
				PacketBody: updateAckBody,
			}
			ackFramePayload, err = packet.Encode(updateAck)
			if err != nil {
				fmt.Println("handlePacket: packet encode error:", err)
				return nil, err
			}
			return ackFramePayload, nil
		}
	case *packet.GetUserDataRequestPacket:
		{
			fmt.Println("fetching user data...")
			getUserDataReq := p.(*packet.GetUserDataRequestPacket)
			nickname := getUserDataReq.PacketBody.Nickname
			sessionId := getUserDataReq.PacketBody.SessionId
			userId, err := strconv.Atoi(sessionId)
			if err != nil {
				fmt.Println(err)
			}

			photoPath, _ := db.PhotoPathById(int64(userId))

			getUserDataAckBody := packet.GetUserDataAckBody{
				nickname,
				photoPath,
			}
			getUserDataAck := &packet.GetUserDataAck{
				ID:         getUserDataReq.ID,
				PacketBody: getUserDataAckBody,
			}
			ackFramePayload, err = packet.Encode(getUserDataAck)
			if err != nil {
				fmt.Println("handlePacket: packet encode error:", err)
				return nil, err
			}
			return ackFramePayload, nil
		}
	case *packet.UpdateNicknameRequest:
		{
			updateNicknameReq := p.(*packet.UpdateNicknameRequest)
			newNickname := updateNicknameReq.PacketBody.Nickname
			sessionId := updateNicknameReq.PacketBody.SessionId

			fmt.Printf("Updating nickname, session id: %s\n", sessionId)
			userId, err := strconv.Atoi(sessionId)
			if err != nil {
				fmt.Println(err)
			}

			err = db.UpdateNicknameById(int64(userId), newNickname)

			updateNicknameAckBody := packet.UpdateNicknameAckBody{
				err == nil,
			}
			updateNicknameAck := &packet.UpdateNicknameAck{
				ID:         updateNicknameReq.ID,
				PacketBody: updateNicknameAckBody,
			}
			ackFramePayload, err = packet.Encode(updateNicknameAck)
			if err != nil {
				fmt.Println("handlePacket: packet encode error:", err)
				return nil, err
			}
			return ackFramePayload, nil
		}
	default:
		return nil, fmt.Errorf("unknown packet type")
	}
}

/* Handlers end */

func main() {
	loadConfigs()
	checkAndStartMySQLServer()
	db.InitDBConnection(configs.Database)

	// Open a socket
	socketAddr := configs.TCPServer.Host + ":" + strconv.Itoa(configs.TCPServer.Port)
	l, err := net.Listen(configs.TCPServer.Protocol, socketAddr)
	if err != nil {
		fmt.Println("TCP server listen socket error:", err)
		return
	}
	fmt.Println("TCP server listen to", socketAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("TCP server socket accept error:", err)
			break
		}

		// start a new goroutine to handle the new connection
		go handleSocketConnection(conn)
	}
}
