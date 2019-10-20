package sql

import (
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"strconv"
)

type Verify struct {
	ChatId   string `gorm:"primary_key"`
	Option   string `gorm:"not null"`
	Deletion string `gorm:"not null"`
	Text     string `gorm:"not null"`
}

func UpdateVerify(chatId int, option string, text string, del string) {
	set := &Verify{}
	st := SESSION.Where(Verify{ChatId: strconv.Itoa(chatId)}).Assign(Verify{Option: option,
		Text: text, Deletion: del}).FirstOrCreate(set)
	err_handler.HandleErr(st.Error)
}

func DelVerify(chatid int) bool {
	filter := &Verify{ChatId: strconv.Itoa(chatid)}

	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func GetVerify(chatid int) *Verify {
	ver := &Verify{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(ver).RowsAffected == 0 {
		return nil
	}
	return ver
}
