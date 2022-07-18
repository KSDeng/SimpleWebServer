package packet

type LoginRequestPacketBody struct {
	Nickname string
	Password string
}

type LoginRequestAckBody struct {
	Nickname  string
	Password  string
	PhotoPath string
	SessionId string
	Result    bool
}

type UpdatePhotoPathPacketBody struct {
	Nickname  string
	PhotoPath string
	SessionId string
}

type UpdatePhotoPathAckBody struct {
	Nickname string
	Succeed  bool
}

type UpdateNicknamePacketBody struct {
	SessionId string
	Nickname  string
}

type UpdateNicknameAckBody struct {
	Succeed bool
}

type GetUserDataRequestBody struct {
	Nickname  string
	SessionId string
}

type GetUserDataAckBody struct {
	Nickname  string
	PhotoPath string
}
