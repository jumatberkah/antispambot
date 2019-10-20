package sql

import (
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"strconv"
)

type Username struct {
	ChatId   string `gorm:"primary_key"`
	Option   string `gorm:"not null"`
	Action   string `gorm:"not null"`
	Deletion string `gorm:"not null"`
	Text     string `gorm:"not null"`
}

func UpdateUsername(chatid int, option string, action string, text string, del string) {
	set := &Username{}
	st := SESSION.Where(Username{ChatId: strconv.Itoa(chatid)}).Assign(Username{Option: option,
		Action: action, Text: text, Deletion: del}).FirstOrCreate(set)
	err_handler.HandleErr(st.Error)
}

func DelUsername(chatid int) bool {
	filter := &Username{ChatId: strconv.Itoa(chatid)}

	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func GetUsername(chatid int) *Username {
	opt := &Username{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(opt).RowsAffected == 0 {
		return nil
	}
	return opt
}
