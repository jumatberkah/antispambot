package err_handler

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/jumatberkah/antispambot/bot/helpers/logger"
	log "github.com/sirupsen/logrus"
)

type CommandCallback func()

func HandleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func FatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func HandleTgErr(b ext.Bot, u *gotgbot.Update, err error) {
	if err != nil {
		err := logger.SendLog(b, u, "error", err.Error())
		HandleErr(err)
	}
}
