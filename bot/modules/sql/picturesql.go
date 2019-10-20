package sql

import (
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"strconv"
)

type Picture struct {
	ChatId   string `gorm:"primary_key"`
	Option   string `gorm:"not null"`
	Action   string `gorm:"not null"`
	Deletion string `gorm:"not null"`
	Text     string `gorm:"not null"`
}

func UpdatePicture(chatid int, option string, action string, text string, del string) {
	set := &Picture{}
	st := SESSION.Where(Picture{ChatId: strconv.Itoa(chatid)}).Assign(Picture{Option: option, Action: action,
		Text: text, Deletion: del}).FirstOrCreate(set)
	err_handler.HandleErr(st.Error)
}

func DelPicture(chatid int) bool {
	filter := &Picture{ChatId: strconv.Itoa(chatid)}

	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func GetPicture(chatid int) *Picture {
	opt := &Picture{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(opt).RowsAffected == 0 {
		return nil
	}
	return opt
}
