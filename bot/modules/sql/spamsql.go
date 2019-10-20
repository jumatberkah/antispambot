package sql

import (
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"strconv"
)

type UserSpam struct {
	UserId int    `gorm:"primary_key"`
	Reason string `gorm:"not null"`
}

type ChatSpam struct {
	ChatId string `gorm:"primary_key"`
	Reason string `gorm:"not null"`
}

type EnforceGban struct {
	ChatId string `gorm:"primary_key"`
	Option string `gorm:"not null"`
}

// User
func UpdateUserSpam(userid int, reason string) {
	user := &UserSpam{}
	st := SESSION.Where(UserSpam{UserId: userid}).Assign(UserSpam{Reason: reason}).FirstOrCreate(user)
	err_handler.HandleErr(st.Error)
}

func DelUserSpam(userid int) bool {
	filter := &UserSpam{UserId: userid}
	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func GetUserSpam(userid int) *UserSpam {
	spam := &UserSpam{UserId: userid}

	if SESSION.First(spam).RowsAffected == 0 {
		return nil
	}
	return spam
}

// Chat
func UpdateChatSpam(chatid int, reason string) {
	chat := &ChatSpam{}
	st := SESSION.Where(ChatSpam{ChatId: strconv.Itoa(chatid)}).Assign(ChatSpam{Reason: reason}).FirstOrCreate(chat)
	err_handler.HandleErr(st.Error)
}

func DelChatSpam(chatid int) bool {
	filter := &ChatSpam{ChatId: strconv.Itoa(chatid)}

	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

// Enforce Gban
func UpdateEnforceGban(chatid int, option string) {
	chat := &EnforceGban{}
	st := SESSION.Where(EnforceGban{ChatId: strconv.Itoa(chatid)}).Assign(EnforceGban{Option: option}).FirstOrCreate(chat)
	err_handler.HandleErr(st.Error)
}

func DelEnforceGban(chatid int) bool {
	filter := &EnforceGban{ChatId: strconv.Itoa(chatid)}

	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

// Get Function
func GetChatSpam(chatid int) *ChatSpam {
	spam := &ChatSpam{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(spam).RowsAffected == 0 {
		return nil
	}
	return spam
}

func GetEnforceGban(chatid int) *EnforceGban {
	spam := &EnforceGban{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(spam).RowsAffected == 0 {
		return nil
	}
	return spam
}

func GetAllSpamUser() []UserSpam {
	var list []UserSpam
	SESSION.Find(&list)
	return list
}

func GetAllSpamChat() []ChatSpam {
	var list []ChatSpam
	SESSION.Find(&list)
	return list
}
