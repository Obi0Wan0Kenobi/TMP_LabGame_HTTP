package main

import (
	"LabRockPaperScissors/db"
	"LabRockPaperScissors/requests"
	"LabRockPaperScissors/structs"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"time"
)

func delete_at_index(slice []structs.ConnectionWS, index int) []structs.ConnectionWS {

	return append(slice[:index], slice[index+1:]...)

}

var Connections []structs.ConnectionWS

func main() {
	rand.Seed(time.Now().UnixNano())
	db.ConnectToDB()

	//fmt.Println(db.FindPlayingRes())
	go db.GameResult()
	go db.CompareEnemys()
	app := fiber.New()

	app.Use(cors.New())

	app.Static("/websocket", "./htmlws")
	app.Static("/", "./html")

	app.Get("/new", requests.NewPlayer)

	app.Get("/find/:uid", requests.FindEnemy)

	app.Get("/choose/:uid", requests.Choseitem)

	app.Get("/check/:uid", requests.Checkwinornot)

	app.Get("/restart/:uid", requests.RestartGame)

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Websocket logic
		check := ""
		for {
			mtype, msg, err := c.ReadMessage()
			//fmt.Println(mtype)
			if err != nil {
				break
			}
			fmt.Println(mtype)
			if mtype == -1 {
				for index, user := range Connections {
					if user.UID == check {
						Connections = delete_at_index(Connections, index)
						db.Connections = delete_at_index(Connections, index)
					}
				}

			}
			var jsonMap structs.MessageFromUser
			err = json.Unmarshal(msg, &jsonMap)
			if err != nil {
				break
			}
			log.Printf("Read: %s", jsonMap)
			if jsonMap.Type == "new" && check == "" {
				user1 := structs.User{0, uuid.New().String(), false, "0", requests.GenToken(), "none", "none"}
				Connections = append(Connections, structs.ConnectionWS{
					Connection: c,
					UID:        user1.Useruid,
				})
				db.Connections = Connections
				check = user1.Useruid
				db.AddUser(user1)

				jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "UserUid": user1.Useruid, "Token": user1.Token})
				err = c.WriteMessage(mtype, jsonstring)
				if err != nil {
					break
				}
			}

			if jsonMap.Type == "choose" {
				user1 := db.FindUserByUID(jsonMap.UID)
				db.Choose(jsonMap.Chose, user1)
				jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Status": true})
				err = c.WriteMessage(mtype, jsonstring)

				if err != nil {
					break
				}

			}

			if jsonMap.Type == "restart" {
				user1 := db.FindUserByUID(jsonMap.UID)
				db.RestartGame(user1)
				jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Restarted": true})
				err = c.WriteMessage(mtype, jsonstring)

				if err != nil {
					break
				}

			}

			if err != nil {
				break
			}
			//if UserId1 != "" {
			//	user := db.FindUserByUID(UserId1)
			//	if user.Playing != false && user.Gamestatus == "0" && Finded == false {
			//		jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Finded": true, "Playingwith": user.Playingwith})
			//		err = c.WriteMessage(mtype, jsonstring)
			//		Finded = true
			//	}
			//}

		}
		//log.Println("Error:", err)
	}))

	app.Listen(":80")
}
