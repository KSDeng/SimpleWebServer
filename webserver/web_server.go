package main

import (
	"encoding/base64"
	"entry-task/app_data_model"
	"entry-task/config"
	"entry-task/network_potocols/frame"
	"entry-task/network_potocols/packet"
	redis_helper "entry-task/webserver/redis"
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var configs config.Configurations

var conn net.Conn
var frameCodec frame.StreamFrameCodec
var requestCounter int

var validPath = regexp.MustCompile("^/(edit|upload)/([\u4e00-\u9fa5a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles("failure.html", "index.html",
	"user_info_display.html", "user_info_edit.html", "user_info_upload.html", "edit_nickname.html"))

/* Utils begin */
func loadConfigs() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("../config")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	err := viper.Unmarshal(&configs)
	if err != nil {
		return err
	}
	return nil
}

func checkAndStartRedisServer() {
	cmd := exec.Command("redis-cli", "ping")
	stdout, _ := cmd.Output()
	output_str := string(stdout)
	if !strings.Contains(output_str, "PONG") {
		cmd_start := exec.Command("redis-server")
		go func() {
			err := cmd_start.Run()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
}

func sendPacketBySocket(p packet.Packet) {
	framePayload, err := packet.Encode(p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("send frame length = %d\n", len(framePayload)+4)
	err = frameCodec.Encode(conn, framePayload)
	if err != nil {
		fmt.Println(err)
	}
}

func handleLoginSucceed(w http.ResponseWriter, r *http.Request, nickname string, photoPath string, sessionId string) {
	fmt.Println("handleLoginSucceed")
	if photoPath == "" {
		data := map[string]interface{}{"Nickname": nickname, "SessionId": sessionId}
		err := templates.ExecuteTemplate(w, "user_info_upload.html", data)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Redirect(w, r, "/display?nickname="+nickname+"&sessionId="+sessionId, http.StatusFound)
	}
}

func handleLoginFailure(w http.ResponseWriter, r *http.Request, nickname string) {
	data := map[string]interface{}{"Nickname": nickname}
	err := templates.ExecuteTemplate(w, "failure.html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func readPhotoBase64(photoPath string) string {
	photoBytes, err := ioutil.ReadFile(photoPath)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(photoBytes)
}

/* Utils end */

/* Handlers begin */
func homePageHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	nickname := r.FormValue("nickname")
	password := r.FormValue("password")

	// send request to tcp server
	requestCounter++
	id := fmt.Sprintf("%08d", requestCounter)
	loginRequestPacketBody := packet.LoginRequestPacketBody{
		nickname,
		password,
	}
	p := &packet.LoginRequestPacket{
		ID:         id,
		PacketBody: loginRequestPacketBody,
	}

	sendPacketBySocket(p)

	// handle reply
	for {
		// handle ack, read from the connection
		ackFramePayload, err := frameCodec.Decode(conn)
		if err != nil {
			fmt.Println(err)
		}

		p, err := packet.Decode(ackFramePayload)

		ack, ok := p.(*packet.LoginRequestAck)

		if !ok {
			fmt.Println("not ack")
		}
		fmt.Println("the result of login request is", ack)
		if ack.ID == id {

			if ack.PacketBody.Result {
				// add data to redis
				redis_helper.SetDataToRedisById(ack.PacketBody.SessionId, app_data_model.UserInfo{
					ack.PacketBody.Nickname,
					ack.PacketBody.Password,
					ack.PacketBody.PhotoPath,
				})

				handleLoginSucceed(w, r, nickname, ack.PacketBody.PhotoPath, ack.PacketBody.SessionId)
			} else {
				handleLoginFailure(w, r, nickname)
			}
			break
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index/", http.StatusFound)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File uploading...")
	m := validPath.FindStringSubmatch(r.URL.Path)
	nickname := m[2]
	sessionId := r.PostFormValue("sessionId")

	// maximum 10M
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")

	if err != nil {
		fmt.Println("Error receiving the file")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Upload file: %+v\n", handler.Filename)

	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	// update database
	requestCounter++
	id := fmt.Sprintf("%08d", requestCounter)

	updatePhotoPathRequestBody := packet.UpdatePhotoPathPacketBody{
		nickname,
		tempFile.Name(),
		sessionId,
	}
	p := &packet.UpdatePhotoPathPacket{
		ID:         id,
		PacketBody: updatePhotoPathRequestBody,
	}

	sendPacketBySocket(p)

	// handle reply
	for {
		// handle ack, read from the connection
		ackFramePayload, err := frameCodec.Decode(conn)
		if err != nil {
			fmt.Println(err)
		}

		p, err := packet.Decode(ackFramePayload)

		ack, ok := p.(*packet.UpdatePhotoPathAck)

		if !ok {
			fmt.Println("not ack")
		}
		fmt.Println("the result of login request is", ack)
		if ack.ID == id {
			if ack.PacketBody.Succeed {
				break
			} else {
				fmt.Println("Update photo path failed")
			}
		}
	}

	// update redis cache
	fmt.Println("Update redis cache photo path to ", tempFile.Name())
	redis_helper.UpdatePhotoPathById(sessionId, tempFile.Name())

	http.Redirect(w, r, "/display?nickname="+nickname+"&sessionId="+sessionId, http.StatusFound)
}

func displayHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nickname := query["nickname"][0]
	sessionId := query["sessionId"][0]

	user_info, res := redis_helper.GetInfoFromRedisById(sessionId)

	if res {
		base64Encoding := readPhotoBase64(user_info.PhotoPath)
		data := map[string]interface{}{"Nickname": nickname, "Image": base64Encoding, "SessionId": sessionId}
		err := templates.ExecuteTemplate(w, "user_info_display.html", data)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// get user data request to TCP server
	requestCounter++
	id := fmt.Sprintf("%08d", requestCounter)

	getUserDataRequestBody := packet.GetUserDataRequestBody{
		nickname,
		sessionId,
	}
	p := &packet.GetUserDataRequestPacket{
		ID:         id,
		PacketBody: getUserDataRequestBody,
	}

	sendPacketBySocket(p)

	// handle reply
	for {
		// handle ack, read from the connection
		ackFramePayload, err := frameCodec.Decode(conn)
		if err != nil {
			fmt.Println(err)
		}

		p, err := packet.Decode(ackFramePayload)

		ack, ok := p.(*packet.GetUserDataAck)

		if !ok {
			fmt.Println("not ack")
		}
		fmt.Println("the result of login request is", ack)
		if ack.ID == id {
			fmt.Println("Update redis cache photo path to ", ack.PacketBody.PhotoPath)
			redis_helper.UpdatePhotoPathById(sessionId, ack.PacketBody.PhotoPath)

			if ack.PacketBody.PhotoPath == "" {
				http.Redirect(w, r, "/upload?nickname="+nickname, http.StatusFound)
				return
			}
			base64Encoding := readPhotoBase64(ack.PacketBody.PhotoPath)
			data := map[string]interface{}{"Nickname": nickname, "Image": base64Encoding, "SessionId": sessionId}
			err := templates.ExecuteTemplate(w, "user_info_display.html", data)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			break
		}
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	m := strings.Split(r.URL.Path, "/")

	nickname := m[2]
	sessionId := m[3]

	data := map[string]interface{}{"Nickname": nickname, "SessionId": sessionId}
	err := templates.ExecuteTemplate(w, "user_info_upload.html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func editNicknameHandler(w http.ResponseWriter, r *http.Request) {
	m := strings.Split(r.URL.Path, "/")

	nickname := m[2]
	sessionId := m[3]

	data := map[string]interface{}{"Nickname": nickname, "SessionId": sessionId}
	err := templates.ExecuteTemplate(w, "edit_nickname.html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveNicknameHandler(w http.ResponseWriter, r *http.Request) {
	newNickname := r.FormValue("nickname")
	sessionId := r.FormValue("sessionId")

	requestCounter++
	id := fmt.Sprintf("%08d", requestCounter)

	updateNicknameReqBody := packet.UpdateNicknamePacketBody{
		sessionId,
		newNickname,
	}
	p := &packet.UpdateNicknameRequest{
		ID:         id,
		PacketBody: updateNicknameReqBody,
	}

	sendPacketBySocket(p)

	// handle reply
	for {
		// handle ack, read from the connection
		ackFramePayload, err := frameCodec.Decode(conn)
		if err != nil {
			fmt.Println("not ack")
		}

		p, err := packet.Decode(ackFramePayload)

		ack, ok := p.(*packet.UpdateNicknameAck)

		if !ok {
			fmt.Println("not ack")
		}

		if ack.ID == id {
			if ack.PacketBody.Succeed {
				redis_helper.UpdateNicknameById(sessionId, newNickname)
				http.Redirect(w, r, "/display?nickname="+newNickname+"&sessionId="+sessionId, http.StatusFound)
			}
			break
		}
	}
}

/* Handlers end */

func main() {
	var err error

	err = loadConfigs()
	if err != nil {
		log.Fatal(err)
	}

	checkAndStartRedisServer()
	redis_helper.InitRedisClient(configs.Redis)

	// start socket
	socketAddr := configs.TCPServer.Host + ":" + strconv.Itoa(configs.TCPServer.Port)
	conn, err = net.Dial(configs.TCPServer.Protocol, socketAddr)

	if err != nil {
		log.Fatal(err)
	}

	frameCodec = frame.NewMyFrameCodec()

	// http requests handler
	http.HandleFunc("/index/", homePageHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/upload/", uploadHandler)
	http.HandleFunc("/display/", displayHandler)
	http.HandleFunc("/logout/", logoutHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/editName/", editNicknameHandler)
	http.HandleFunc("/saveNickname/", saveNicknameHandler)

	// Listen to http request from client
	addr := configs.WebServer.Host + ":" + strconv.Itoa(configs.WebServer.Port)
	fmt.Println("Web server listening on", addr)

	go log.Fatal(http.ListenAndServe(addr, nil))

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)
}
