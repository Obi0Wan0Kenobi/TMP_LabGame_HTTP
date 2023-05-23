package structs

import "github.com/gofiber/websocket/v2"

type User struct {
	Id          int
	Useruid     string
	Playing     bool
	Chosen      string
	Token       string
	Playingwith string
	Gamestatus  string
}

type MessageFromUser struct {
	Type  string `json:"type"`
	Token string `json:"token"`
	UID   string `json:"uid"`
	Chose string `json:"chose"`
}

type ConnectionWS struct {
	Connection *websocket.Conn `json:"type"`
	UID        string          `json:"uid"`
}
