package modules

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"regexp"
	"strconv"
	"strings"
)

func panel(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat.Type == "supergroup" {
		if chat_status.IsUserAdmin(chat, user.Id) == true {
			teks, _, kn := function.MainMenu(chat.Id)
			reply := b.NewSendableMessage(chat.Id, teks)
			reply.ReplyMarkup = &ext.InlineKeyboardMarkup{&kn}
			reply.ParseMode = parsemode.Html
			reply.ReplyToMessageId = msg.MessageId
			_, err = reply.Send()
			return err
		}
	}
	return nil
}

func backquery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				teks, _, kn := function.MainMenu(chat.Id)
				_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId, teks, parsemode.Html,
					&ext.InlineKeyboardMarkup{&kn})
				return err
			}
		} else {
			_, err = b.AnswerCallbackQuery(msg.Id)
			return err
		}
	}
	return gotgbot.ContinueGroups{}
}

func closequery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				_, err = msg.Message.Delete()
				return err
			}
		} else {
			_, err = b.AnswerCallbackQuery(msg.Id)
			return err
		}
	}
	return gotgbot.ContinueGroups{}
}

func settingquery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				if msg.Data == "mk_utama" {
					teks, _, kn := function.MainControlMenu(chat.Id)
					_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId,
						teks, "HTML", &ext.InlineKeyboardMarkup{&kn})
					return err
				} else if msg.Data == "mk_reset" {
					sql.UpdatePicture(chat.Id, "true", "mute", "-", "true")
					sql.UpdateUsername(chat.Id, "true", "mute", "-", "true")
					sql.UpdateEnforceGban(chat.Id, "true")
					sql.UpdateVerify(chat.Id, "true", "-", "true")
					sql.UpdateSetting(chat.Id, "5m", "true")

					err = updateusercontrol(b, u)
					return err
				} else if msg.Data == "mk_spam" {
					teks, _, kn := function.MainSpamMenu(chat.Id)
					_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId,
						teks, "HTML", &ext.InlineKeyboardMarkup{&kn})
					return err
				} else {
					_, err = b.AnswerCallbackQuery(msg.Id)
				}
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func usercontrolquery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			if chat_status.IsUserAdmin(chat, user.Id) == true {
				// Grab Data From DB
				username := sql.GetUsername(chat.Id)
				fotoprofil := sql.GetPicture(chat.Id)
				waktu := sql.GetSetting(chat.Id)
				ver := sql.GetVerify(chat.Id)

				// Separating Queries
				z, _ := regexp.MatchString("^m[cdef]_del$", msg.Data)
				a, _ := regexp.MatchString("^mc_(kick|ban|mute)$", msg.Data)
				f, _ := regexp.MatchString("^md_(kick|ban|mute)$", msg.Data)
				g, _ := regexp.MatchString("^m[cdeo]_toggle$", msg.Data)
				d, _ := regexp.MatchString("^mf_(plus|minus|duration|waktu)$", msg.Data)

				// Username Control Panel
				if a == true {
					sql.UpdateUsername(chat.Id, username.Option, strings.Split(msg.Data, "mc_")[1], "-", username.Deletion)
					err = updateusercontrol(b, u)
					return err
				} else if f == true {
					// Profile Photo Control Panel
					sql.UpdatePicture(chat.Id, fotoprofil.Option, strings.Split(msg.Data, "md_")[1], "-", fotoprofil.Deletion)
					err = updateusercontrol(b, u)
					return err
				} else if d == true {
					// Time Control Panel
					if strings.Split(msg.Data, "mf_")[1] == "duration" {
						lastLetter := waktu.Time[len(waktu.Time)-1:]
						lastLetter = strings.ToLower(lastLetter)
						re := regexp.MustCompile(`[mhd]`)

						t := waktu.Time[:len(waktu.Time)-1]
						_, err := strconv.Atoi(t)
						if err != nil {
							_, err := b.AnswerCallbackQueryText(msg.Id,
								"❌ Invalid time amount specified.", true)
							return err
						}

						if lastLetter == "m" {
							sql.UpdateSetting(chat.Id, fmt.Sprintf("%vh", re.Split(waktu.Time, -1)[0]), waktu.Deletion)
						} else if lastLetter == "h" {
							sql.UpdateSetting(chat.Id, fmt.Sprintf("%vd", re.Split(waktu.Time, -1)[0]), waktu.Deletion)
						} else if lastLetter == "d" {
							sql.UpdateSetting(chat.Id, fmt.Sprintf("%vm", re.Split(waktu.Time, -1)[0]), waktu.Deletion)
						}

						err = updateusercontrol(b, u)
						return err

					} else if strings.Split(msg.Data, "mf_")[1] == "plus" {
						lastLetter := waktu.Time[len(waktu.Time)-1:]
						lastLetter = strings.ToLower(lastLetter)

						t := waktu.Time[:len(waktu.Time)-1]
						j, err := strconv.Atoi(t)
						if err != nil {
							_, err := b.AnswerCallbackQueryText(msg.Id,
								"❌ Invalid time amount specified.", true)
							return err
						} else {
							j++
						}

						sql.UpdateSetting(chat.Id, fmt.Sprintf("%v%v", j, lastLetter), waktu.Deletion)
						err = updateusercontrol(b, u)
						return err
					} else if strings.Split(msg.Data, "mf_")[1] == "minus" {
						lastLetter := waktu.Time[len(waktu.Time)-1:]
						lastLetter = strings.ToLower(lastLetter)
						if strings.ContainsAny(lastLetter, "m & d & h") {
							t := waktu.Time[:len(waktu.Time)-1]
							j, err := strconv.Atoi(t)
							if err != nil {
								_, err := b.AnswerCallbackQueryText(msg.Id,
									"❌ Invalid time amount specified.", true)
								return err
							} else {
								j--
							}

							sql.UpdateSetting(chat.Id, fmt.Sprintf("%v%v", j, lastLetter), waktu.Deletion)
							err = updateusercontrol(b, u)
							return err
						}
					} else if strings.Split(msg.Data, "mf_")[1] == "waktu" {
						_, err := b.AnswerCallbackQueryText(msg.Id,
							"🔄 Mengatur tenggat waktu untuk semua aksi.", true)
						return err
					}
				} else if g == true {
					// On/Off Toggles
					if strings.Split(msg.Data, "_toggle")[0] == "mc" {
						if username.Option == "true" {
							sql.UpdateUsername(chat.Id, "false", username.Action, "-", username.Deletion)
						} else {
							sql.UpdateUsername(chat.Id, "true", username.Action, "-", username.Deletion)
						}
					} else if strings.Split(msg.Data, "_toggle")[0] == "md" {
						if fotoprofil.Option == "true" {
							sql.UpdatePicture(chat.Id, "false", username.Action, "-", username.Deletion)
						} else {
							sql.UpdatePicture(chat.Id, "true", username.Action, "-", username.Deletion)
						}
					} else if strings.Split(msg.Data, "_toggle")[0] == "me" {
						if ver.Option == "true" {
							sql.UpdateVerify(chat.Id, "false", "-", ver.Deletion)
						} else {
							sql.UpdateVerify(chat.Id, "true", "-", ver.Deletion)
						}
					} else if strings.Split(msg.Data, "_toggle")[0] == "mo" {
						if fotoprofil.Option == "true" {
							sql.UpdateEnforceGban(chat.Id, "false")
						} else {
							sql.UpdateEnforceGban(chat.Id, "true")
						}
					}

					err = updateusercontrol(b, u)
					return err
				} else if z == true {
					// On/Off Deletion
					if strings.Split(msg.Data, "_del")[0] == "mc" {
						if username.Deletion == "true" {
							sql.UpdateUsername(chat.Id, username.Option, username.Action, "-", "false")
						} else {
							sql.UpdateUsername(chat.Id, username.Option, username.Action, "-", "true")
						}
					} else if strings.Split(msg.Data, "_del")[0] == "md" {
						if fotoprofil.Deletion == "true" {
							sql.UpdatePicture(chat.Id, fotoprofil.Option, fotoprofil.Action, "-", "false")
						} else {
							sql.UpdatePicture(chat.Id, fotoprofil.Option, fotoprofil.Action, "-", "true")
						}
					} else if strings.Split(msg.Data, "_del")[0] == "me" {
						if ver.Deletion == "true" {
							sql.UpdateVerify(chat.Id, ver.Option, "-", "false")
						} else {
							sql.UpdateVerify(chat.Id, ver.Option, "-", "true")
						}
					} else if strings.Split(msg.Data, "_del")[0] == "mf" {
						if waktu.Deletion == "true" {
							sql.UpdateSetting(chat.Id, waktu.Time, "false")
						} else {
							sql.UpdateSetting(chat.Id, waktu.Time, "true")
						}
					}

					err = updateusercontrol(b, u)
					return err
				}
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func updateusercontrol(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	chat := msg.Message.Chat

	// Strings
	opsisama := "Bad Request: message is not modified: specified new message content and " +
		"reply markup are exactly the same as a current " +
		"content and reply markup of the message"

	// Main Action
	teks, _, kn := function.MainControlMenu(chat.Id)
	_, err = b.EditMessageTextMarkup(chat.Id, msg.Message.MessageId,
		teks, "HTML", &ext.InlineKeyboardMarkup{&kn})
	if err != nil {
		if err.Error() == opsisama {
			_, err := b.AnswerCallbackQueryText(msg.Id, "❌ Anda memilih pilihan yang sama",
				false)
			return err
		} else {
			_, err := b.AnswerCallbackQueryText(msg.Id, "🔄 Silahkan Coba Lagi",
				true)
			return err
		}
	} else {
		_, err := b.AnswerCallbackQuery(msg.Id)
		return err
	}
}

func LoadSettingPanel(u *gotgbot.Updater) {
	u.Dispatcher.AddHandler(handlers.NewCommand("setting", panel))
	u.Dispatcher.AddHandler(handlers.NewCallback(
		regexp.MustCompile("^m[cdefgo]_(toggle|kick|ban|mute|reset|plus|minus|duration|waktu|del)").String(),
		usercontrolquery))
	u.Dispatcher.AddHandler(handlers.NewCallback(regexp.MustCompile("^mk_").String(),
		settingquery))
	u.Dispatcher.AddHandler(handlers.NewCallback(regexp.MustCompile("^close").String(),
		closequery))
	u.Dispatcher.AddHandler(handlers.NewCallback(regexp.MustCompile("^back").String(),
		backquery))

}
