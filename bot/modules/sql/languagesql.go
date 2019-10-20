package sql

import (
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"strconv"
)

type Lang struct {
	ChatId string `gorm:"primary_key"`
	Lang   string `gorm:"not null"`
}

func UpdateLang(chatid int, lang string) {
	set := &Lang{}
	st := SESSION.Where(Lang{ChatId: strconv.Itoa(chatid)}).Assign(Lang{Lang: lang}).FirstOrCreate(set)
	err_handler.HandleErr(st.Error)
}

func DelLang(chatid int) bool {
	filter := &Lang{ChatId: strconv.Itoa(chatid)}

	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func GetLang(chatid int) *Lang {
	ver := &Lang{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(ver).RowsAffected == 0 {
		return nil
	}
	return ver
}
