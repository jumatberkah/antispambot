package modules

import (
	"fmt"
	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/go-redis/redis"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
)


func GetString (chat_id int, val string) string {
	var err error
	lang, err := sql.REDIS.Get(fmt.Sprintf("lang_%s", chat_id)).Result()
	if err == redis.Nil{
		lang = "en-GB"
	}
	err_handler.HandleErr(err)
	return goloc.Trnl(lang, val)
}

func GetStringf (chat_id int, val string, args map[string]string) string {
	var err error
	lang, err := sql.REDIS.Get(fmt.Sprintf("lang_%s", chat_id)).Result()
	if err == redis.Nil{
		lang = "en-GB"
	}
	err_handler.HandleErr(err)
	return goloc.Trnlf(lang, val, args)
}


func setlang (b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat_status.IsUserAdmin(chat, user.Id, nil) {
		if goloc.IsLangSupported(args[0]) {
			_, err := sql.REDIS.Set(fmt.Sprintf("lang_%s", chat.Id), args[0], 0).Result()
			if err != nil {
				sql.UpdateLang(chat.Id, args[0])
				_, err = msg.ReplyText("Language has been changed.")
				return err
			} else {
				_, err = msg.ReplyText("Language has been changed.")
				return err
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func LoadLang (u *gotgbot.Updater) {
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("setlang", setlang))
}
