package sql

import (
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"strings"
)

type User struct {
	UserId   int `gorm:"primary_key"`
	UserName string
	Name     string `gorm:"not null"`
}

type Chat struct {
	ChatId    string `gorm:"primary_key"`
	ChatTitle string `gorm:"not null"`
	ChatType  string `gorm:"not null"`
	ChatLink  string
}

func UpdateUser(userid int, username string, name string) {
	username = strings.ToLower(username)
	user := &User{}
	st := SESSION.Where(User{UserId: userid}).Assign(User{UserName: username, Name: name}).FirstOrCreate(user)
	err_handler.HandleErr(st.Error)
}

func UpdateChat(chatid string, chattitle string, chattype string, clink string) {
	if chatid == "nil" || chattitle == "nil" {
		return
	}

	chat := &Chat{}
	st := SESSION.Where(Chat{ChatId: chatid}).Assign(Chat{ChatTitle: chattitle, ChatType: chattype,
		ChatLink: clink}).FirstOrCreate(chat)
	err_handler.HandleErr(st.Error)
}

func DelUser(userid int) bool {
	filter := &User{UserId: userid}

	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func DelChat(chatid string) bool {
	filter := &Chat{ChatId: chatid}

	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func GetUserIdByName(username string) *User {
	username = strings.ToLower(username)
	user := new(User)
	SESSION.Where("user_name = ?", username).First(user)
	return user
}

func GetAllChat() []Chat {
	var list []Chat
	SESSION.Find(&list)
	return list
}

func GetAllUser() []User {
	var list []User
	SESSION.Find(&list)
	return list
}
