package modules

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"github.com/jumatberkah/antispambot/bot/helpers/logger"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"strconv"
	"strings"
)

func gban(bot ext.Bot, u *gotgbot.Update, args []string) error {
	var err error
	msg := u.EffectiveMessage

	if chat_status.IsOwner(msg.From.Id) == true {
		userid, reason := extraction.ExtractUserAndText(msg, args)

		if userid == 0 {
			_, err = msg.ReplyHTML(GetString(msg.Chat.Id, "modules/admins.go:26"))
			return err
		}

		if reason == "" {
			reason = "None"
		}

		ban := sql.GetUserSpam(userid)
		group := sql.GetAllChat
		banerr := []string{"Bad Request: USER_ID_INVALID", "Bad Request: USER_NOT_PARTICIPANT" +
			"Bad Request: chat member status can't be changed in private chats"}

		if ban != nil {
			if ban.Reason == reason {
				_, err = msg.ReplyHTML(GetString(msg.Chat.Id, "modules/admins.go:41"))
				return gotgbot.EndGroups{}
			} else {
				go sql.UpdateUserSpam(userid, reason)

				_, err = msg.ReplyHTML(GetStringf(msg.Chat.Id, "modules/admins.go:46",
					map[string]string{"1": strconv.Itoa(userid), "2": reason}))
				err_handler.HandleErr(err)
				err = logger.SendBanLog(bot, userid, reason, u)
				return err
			}
		} else {
			go sql.UpdateUserSpam(userid, reason)

			_, err = msg.ReplyHTML(GetStringf(msg.Chat.Id, "modules/admins.go:55", map[string]string{"1": strconv.Itoa(userid)}))
			err_handler.HandleErr(err)

			go func() {
				for _, a := range group() {
					cid, _ := strconv.Atoi(a.ChatId)
					if sql.GetEnforceGban(cid) != nil && sql.GetEnforceGban(cid).Option == "true" {
						_, err = bot.KickChatMember(cid, userid)
						if err != nil {
							if function.Contains(banerr, err.Error()) == true {
								return
							} else if err.Error() == "Forbidden: bot was kicked from the supergroup chat" {
								sql.DelChat(a.ChatId)
								return
							}
						}
					}
				}
			}()

			_, err = msg.ReplyHTML(GetStringf(msg.Chat.Id, "modules/admins.go:75", map[string]string{"1": strconv.Itoa(userid), "2": reason}))
			err_handler.HandleErr(err)
			err = logger.SendBanLog(bot, userid, reason, u)
			return err
		}
	} else {
		_, err := msg.ReplyHTML(GetString(msg.Chat.Id, "modules/admins.go:81"))
		return err
	}
}

func ungban(bot ext.Bot, u *gotgbot.Update, args []string) error {
	var err error
	msg := u.EffectiveMessage

	if chat_status.IsOwner(msg.From.Id) == true {
		userid := extraction.ExtractUser(msg, args)

		if userid == 0 {
			_, err = msg.ReplyHTML("👤 <i>No user was designated.</i>")
			return err
		}

		group := sql.GetAllChat
		banerr := []string{"Bad Request: USER_ID_INVALID", "Bad Request: USER_NOT_PARTICIPANT" +
			"Bad Request: chat member status can't be changed in private chats"}
		ban := sql.GetUserSpam(userid)

		if ban != nil {
			_, err := msg.ReplyHTMLf("⚠ <b>Beginning unglobal ban procedure(s) of</b>\n"+
				"👤 <b>User ID:</b> <code>%v</code>", userid)
			err_handler.HandleErr(err)

			go func() {
				sql.DelUserSpam(userid)
				for _, a := range group() {
					cid, _ := strconv.Atoi(a.ChatId)
					_, err = bot.UnbanChatMember(cid, userid)
					if err != nil {
						if function.Contains(banerr, err.Error()) == true {
							return
						} else if err.Error() == "Forbidden: bot was kicked from the supergroup chat" {
							go sql.DelChat(a.ChatId)
							return
						}
					}
				}
			}()

			_, err = msg.ReplyHTMLf("<b>New Ungban</b>\n"+
				"<b>👤 User ID:</b> <code>%v</code>", userid)
			return err
		} else {
			_, err = msg.ReplyHTMLf("⚠ <b>User has not been banned.</b>")
			return err
		}

	} else {
		_, err := msg.ReplyHTML("❌ <i>You Are Not Sudoers</i>")
		return err
	}
}

func stats(bot ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	teks := fmt.Sprintf("<b>Statistik</b>\n"+
		"Jml Pengguna: %v\nJml Chat: %v\nJml Spammer: %v", len(sql.GetAllUser()),
		len(sql.GetAllChat()), len(sql.GetAllSpamUser()))

	_, err := msg.ReplyHTML(teks)
	return err
}

func broadcast(bot ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.EffectiveMessage
	group := sql.GetAllChat

	if chat_status.IsOwner(msg.From.Id) == true {
		errnum := 0

		for _, a := range group() {
			cid, _ := strconv.Atoi(a.ChatId)
			_, err = bot.SendMessageHTML(cid, strings.Split(msg.Text, "/broadcast")[1])
			if err != nil {
				if err.Error() == "Forbidden: bot was kicked from the supergroup chat" {
					go sql.DelChat(a.ChatId)
					errnum++
				} else {
					err_handler.HandleErr(err)
					errnum++
				}
			}
		}

		_, err = msg.ReplyHTMLf("<b>Message Has Been Broadcasted</b>, <code>%v</code> <b>Has Failed</b>\n", errnum)
		return err
	} else {
		_, err := msg.ReplyHTML("❌ <i>You Are Not Sudoers</i>")
		return err
	}
}

func LoadAdmins(u *gotgbot.Updater) {
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("gban", gban))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("ungban", ungban))
	u.Dispatcher.AddHandler(handlers.NewCommand("stats", stats))
	u.Dispatcher.AddHandler(handlers.NewCommand("broadcast", broadcast))
}
