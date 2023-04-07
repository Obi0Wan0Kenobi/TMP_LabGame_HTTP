package requests

import (
	"LabRockPaperScissors/db"
	"LabRockPaperScissors/structs"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func GenToken() string {
	n := 40
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NewPlayer(c *fiber.Ctx) error {
	user1 := structs.User{0, uuid.New().String(), false, "0", GenToken(), "none", "none"}
	db.AddUser(user1)
	return c.JSON(map[string]interface{}{"AnswerId": 1, "UserUid": user1.Useruid, "Token": user1.Token})

}

func FindEnemy(c *fiber.Ctx) error {
	user := db.FindUserByUID(c.Params("uid"))
	token := c.Query("token")
	if user.Token != token {
		return c.JSON(map[string]interface{}{"AnswerId": 2, "Finded": false, "Playingwith": ""})
	}
	fmt.Println(user.Playing, user.Gamestatus)
	if user.Playing != false && user.Gamestatus == "0" {
		return c.JSON(map[string]interface{}{"AnswerId": 1, "Finded": true, "Playingwith": user.Playingwith})
	}

	return c.JSON(map[string]interface{}{"AnswerId": 1, "Finded": false, "Playingwith": ""})
}

func Choseitem(c *fiber.Ctx) error {
	user := db.FindUserByUID(c.Params("uid"))
	token := c.Query("token")
	chose := c.Query("chose")

	if !(chose == "1" || chose == "2" || chose == "3") {
		return c.JSON(map[string]interface{}{"AnswerId": 3, "Status": false})
	}
	if user.Token != token {
		return c.JSON(map[string]interface{}{"AnswerId": 2, "Status": false})
	}
	db.Choose(chose, user)

	return c.JSON(map[string]interface{}{"AnswerId": 1, "Status": true})
}

func Checkwinornot(c *fiber.Ctx) error {
	user := db.FindUserByUID(c.Params("uid"))
	token := c.Query("token")

	if user.Token != token {
		return c.JSON(map[string]interface{}{"AnswerId": 2, "Wait": true, "Windata": ""})
	}
	if user.Gamestatus != "0" {
		return c.JSON(map[string]interface{}{"AnswerId": 1, "Wait": false, "Windata": user.Gamestatus})
	}

	return c.JSON(map[string]interface{}{"AnswerId": 1, "Wait": true, "Windata": ""})
}

func RestartGame(c *fiber.Ctx) error {
	user := db.FindUserByUID(c.Params("uid"))
	token := c.Query("token")

	if user.Token != token {
		return c.JSON(map[string]interface{}{"AnswerId": 2, "Restarted": false})
	}

	db.RestartGame(user)

	return c.JSON(map[string]interface{}{"AnswerId": 1, "Restarted": true})
}
