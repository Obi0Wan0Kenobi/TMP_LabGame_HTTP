package main

import (
	"LabRockPaperScissors/db"
	"LabRockPaperScissors/requests"
	"LabRockPaperScissors/structs"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	db.ConnectToDB()

	//fmt.Println(db.FindPlayingRes())
	go db.GameResult()
	go db.CompareEnemys()
	app := fiber.New()

	app.Use(cors.New())

	app.Static("/", "./html")

	app.Get("/new", requests.NewPlayer)

	app.Get("/find/:uid", requests.FindEnemy)

	app.Get("/choose/:uid", requests.Choseitem)

	app.Get("/check/:uid", requests.Checkwinornot)

	app.Get("/restart/:uid", requests.RestartGame)

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Websocket logic
		for {
			mtype, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			var jsonMap structs.MessageFromUser
			err = json.Unmarshal(msg, &jsonMap)
			if err != nil {
				break
			}
			log.Printf("Read: %s", jsonMap)
			if jsonMap.Type == "new" {
				user1 := structs.User{0, uuid.New().String(), false, "0", requests.GenToken(), "none", "none"}
				db.AddUser(user1)
				jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "UserUid": user1.Useruid, "Token": user1.Token})
				err = c.WriteMessage(mtype, jsonstring)
				if err != nil {
					break
				}
			}
			//err = c.WriteMessage(mtype, msg)
			if err != nil {
				break
			}
		}
		//log.Println("Error:", err)
	}))

	app.Listen(":80")
}
