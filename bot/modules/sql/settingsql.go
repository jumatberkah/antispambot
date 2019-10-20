package sql

import (
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"strconv"
)

type Setting struct {
	ChatId   string `gorm:"primary_key"`
	Time     string `gorm:"not null"`
	Deletion string `gorm:"not null"`
}

func UpdateSetting(chatid int, time string, delete string) {
	set := &Setting{}
	st := SESSION.Where(Setting{ChatId: strconv.Itoa(chatid)}).Assign(Setting{Time: time, Deletion: delete}).FirstOrCreate(set)
	err_handler.HandleErr(st.Error)
}

func DelSetting(ChatId int) bool {
	filter := &Setting{ChatId: strconv.Itoa(ChatId)}

	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func GetSetting(chatid int) *Setting {
	tim := &Setting{ChatId: strconv.Itoa(chatid)}

	if SESSION.First(tim).RowsAffected == 0 {
		return nil
	}
	return tim
}
