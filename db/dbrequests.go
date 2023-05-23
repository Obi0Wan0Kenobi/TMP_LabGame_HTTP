package db

import (
	"LabRockPaperScissors/structs"
	"encoding/json"
	"fmt"
	"time"
)

var Connections []structs.ConnectionWS

func RemovefromMas(list []string, item string) []string {
	for i, v := range list {
		if v == item {
			copy(list[i:], list[i+1:])
			//list[len(list)-1] = "" // обнуляем "хвост"
			list = list[:len(list)-1]
		}
	}
	return list
}

func AddUser(user structs.User) {
	_, _ = DB.Exec("INSERT INTO players (useruid,token) values ($1,$2)", user.Useruid, user.Token)
}

func FindUserByUID(uid string) structs.User {
	row := DB.QueryRow("SELECT * FROM players Where useruid = $1", uid)
	user := structs.User{}
	_ = row.Scan(&user.Id, &user.Useruid, &user.Playing, &user.Chosen, &user.Token, &user.Playingwith, &user.Gamestatus)

	return user
}

func FindNotPlaying() structs.User {
	row := DB.QueryRow("SELECT * FROM players Where playing = 0")
	enemy1 := structs.User{}
	_ = row.Scan(&enemy1.Id, &enemy1.Useruid, &enemy1.Playing, &enemy1.Chosen, &enemy1.Token, &enemy1.Playingwith, &enemy1.Gamestatus)

	return enemy1
}

func FindNotPlaying1(id string) structs.User {
	row := DB.QueryRow("SELECT * FROM players Where playing = 0 and useruid!=$1", id)
	enemy1 := structs.User{}
	_ = row.Scan(&enemy1.Id, &enemy1.Useruid, &enemy1.Playing, &enemy1.Chosen, &enemy1.Token, &enemy1.Playingwith, &enemy1.Gamestatus)

	return enemy1
}

func CompareEnemys() {
	for true {
		enemy1 := FindNotPlaying()
		enemy2 := FindNotPlaying1(enemy1.Useruid)
		if enemy2.Useruid != "" {
			fmt.Println("finded")
			_, _ = DB.Exec("UPDATE players SET playing=1,playingwith=$1 where useruid=$2", enemy2.Useruid, enemy1.Useruid)
			_, _ = DB.Exec("UPDATE players SET playing=1,playingwith=$1 where useruid=$2", enemy1.Useruid, enemy2.Useruid)
			for _, item := range Connections {
				if item.UID == enemy1.Useruid {
					jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Finded": true, "Playingwith": enemy2.Useruid})
					item.Connection.WriteMessage(1, jsonstring)
				}
				if item.UID == enemy2.Useruid {
					jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Finded": true, "Playingwith": enemy1.Useruid})
					item.Connection.WriteMessage(1, jsonstring)
				}
			}

		} else {
			//fmt.Println(enemy2.Useruid)
		}
		time.Sleep(1 * time.Second)
	}

}

func FindPlayingRes() []string {
	rows, _ := DB.Query("SELECT * FROM players Where playing = 1  and gamestatus='0'")
	var mas []string
	for rows.Next() {
		var p structs.User
		_ = rows.Scan(&p.Id, &p.Useruid, &p.Playing, &p.Chosen, &p.Token, &p.Playingwith, &p.Gamestatus)
		if p.Useruid != "" {
			mas = append(mas, p.Useruid)
		}
	}
	return mas
}

func Choose(chose string, user structs.User) {
	_, _ = DB.Exec("UPDATE players SET chosen=$1 where useruid=$2 and chosen='0'", chose, user.Useruid)
}

func RestartGame(user structs.User) {
	_, _ = DB.Exec("UPDATE players SET playing=0,chosen='0',playingwith='0',gamestatus='0' where useruid=$1 and gamestatus!='0'", user.Useruid)
}

func Win1_Lose2(user1 structs.User, user2 structs.User) {
	_, _ = DB.Exec("UPDATE players SET gamestatus=$1 where useruid=$2", "Вы выиграли!", user1.Useruid)
	_, _ = DB.Exec("UPDATE players SET gamestatus=$1 where useruid=$2", "Вы проиграли!", user2.Useruid)

	for _, item := range Connections {
		if item.UID == user1.Useruid {

			jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Wait": false, "Windata": "Вы выиграли!"})
			item.Connection.WriteMessage(1, jsonstring)
		}
		if item.UID == user2.Useruid {
			jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Wait": false, "Windata": "Вы проиграли!"})
			item.Connection.WriteMessage(1, jsonstring)
		}
	}

}
func Win2_Lose1(user1 structs.User, user2 structs.User) {
	_, _ = DB.Exec("UPDATE players SET gamestatus=$1 where useruid=$2", "Вы проиграли!", user1.Useruid)
	_, _ = DB.Exec("UPDATE players SET gamestatus=$1 where useruid=$2", "Вы выиграли!", user2.Useruid)

	for _, item := range Connections {
		if item.UID == user1.Useruid {

			jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Wait": false, "Windata": "Вы проиграли!"})
			item.Connection.WriteMessage(1, jsonstring)
		}
		if item.UID == user2.Useruid {
			jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Wait": false, "Windata": "Вы выиграли!"})
			item.Connection.WriteMessage(1, jsonstring)
		}
	}
}
func Draw(user1 structs.User, user2 structs.User) {
	_, _ = DB.Exec("UPDATE players SET gamestatus=$1 where useruid=$2", "Ничья", user1.Useruid)
	_, _ = DB.Exec("UPDATE players SET gamestatus=$1 where useruid=$2", "Ничья", user2.Useruid)
	for _, item := range Connections {
		if item.UID == user1.Useruid {

			jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Wait": false, "Windata": "Ничья"})
			item.Connection.WriteMessage(1, jsonstring)
		}
		if item.UID == user2.Useruid {
			jsonstring, _ := json.Marshal(map[string]interface{}{"AnswerId": 1, "Wait": false, "Windata": "Ничья"})
			item.Connection.WriteMessage(1, jsonstring)
		}
	}
}

func GameResult() {
	for true {
		enemys := FindPlayingRes()
		for len(enemys) > 0 {
			enemy1 := FindUserByUID(enemys[0])
			enemy2 := FindUserByUID(enemy1.Playingwith)
			enemys = RemovefromMas(enemys, enemy1.Useruid)
			enemys = RemovefromMas(enemys, enemy2.Useruid)
			//1- камень 2- ножницы 3- бумага
			if enemy1.Chosen != "0" && enemy2.Chosen != "0" {
				if enemy1.Chosen == enemy2.Chosen {
					Draw(enemy1, enemy2)
				} else if (enemy1.Chosen == "1" && enemy2.Chosen == "2") || (enemy1.Chosen == "2" && enemy2.Chosen == "3") || (enemy1.Chosen == "3" && enemy2.Chosen == "1") {
					Win1_Lose2(enemy1, enemy2)
				} else {
					Win2_Lose1(enemy1, enemy2)
				}
			}

		}

		time.Sleep(1 * time.Second)

	}
}
