package modules

import (
	"encoding/json"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/extraction"
	"github.com/jumatberkah/antispambot/bot/helpers/logger"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type spammer struct {
	Status bool `json:"ok"`
}

func username(b ext.Bot, u *gotgbot.Update) error {
	var err error
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	db := sql.GetUsername(chat.Id)

	if db.Option != "true" {
		return gotgbot.EndGroups{}
	}
	if chat_status.IsUserAdmin(chat, msg.From.Id) == true {
		return gotgbot.EndGroups{}
	}

	if msg != nil {
		if chat.Type == "supergroup" {
			if user.Username == "" {
				bantime := extraction.ExtractTime(b, msg, sql.GetSetting(chat.Id).Time)
				replytext := GetStringf(msg.Chat.Id, "modules/listener.go:45",
					map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": db.Action,
						"4": strconv.Itoa(user.Id)})

				kb := make([][]ext.InlineKeyboardButton, 1)
				kb[0] = make([]ext.InlineKeyboardButton, 1)
				kb[0][0] = ext.InlineKeyboardButton{Text: GetString(chat.Id, "modules/listener.go:51"),
					CallbackData: fmt.Sprintf("umute_%v", user.Id)}

				kbk := make([][]ext.InlineKeyboardButton, 1)
				kbk[0] = make([]ext.InlineKeyboardButton, 1)
				kbk[0][0] = ext.InlineKeyboardButton{Text: GetString(chat.Id, "modules/listener.go:56"),
					CallbackData: fmt.Sprintf("uba_%v", user.Id)}

				// Implementing The Chosen Action To The Target
				if db.Action == "mute" {
					restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
					restrictSend.UntilDate = bantime
					_, err = restrictSend.Send()
					if err != nil {
						if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
							_, err = b.SendMessage(chat.Id, err.Error())
							return err
						}
					}
					// Send Message
					reply := b.NewSendableMessage(chat.Id, replytext)
					reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
					reply.ParseMode = parsemode.Html
					reply.ReplyToMessageId = msg.MessageId
					_, err = reply.Send()
					if err != nil {
						if err.Error() == "Bad Request: reply message not found" {
							reply := b.NewSendableMessage(chat.Id, replytext)
							reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
							reply.ParseMode = parsemode.Html
							_, err = reply.Send()
							return err
						} else {
							err_handler.HandleErr(err)
						}
					}
				} else if db.Action == "kick" {
					restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
					_, err = restrictSend.Send()
					if err != nil {
						if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
							_, err = b.SendMessage(chat.Id, err.Error())
							return err
						}
					}
					_, err = b.UnbanChatMember(chat.Id, user.Id)
					if err != nil {
						if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
							_, err = b.SendMessage(chat.Id, err.Error())
							return err
						}
					}
					// Send Message
					reply := b.NewSendableMessage(chat.Id, replytext)
					reply.ParseMode = parsemode.Html
					reply.ReplyToMessageId = msg.MessageId
					_, err = reply.Send()
					return err
				} else if db.Action == "ban" {
					restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
					restrictSend.UntilDate = -1
					_, err = restrictSend.Send()
					if err != nil {
						if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
							_, err = b.SendMessage(chat.Id, err.Error())
							return err
						}
					}
					// Send Message
					reply := b.NewSendableMessage(chat.Id, replytext)
					reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kbk}
					reply.ParseMode = parsemode.Html
					reply.ReplyToMessageId = msg.MessageId
					_, err = reply.Send()
					if err != nil {
						if err.Error() == "Bad Request: reply message not found" {
							reply := b.NewSendableMessage(chat.Id, replytext)
							reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kbk}
							reply.ParseMode = parsemode.Html
							_, err = reply.Send()
							return err
						} else {
							err_handler.HandleErr(err)
						}
					}
				}
				// Delete His/Her Message(s)
				if db.Deletion == "true" {
					_, err = msg.Delete()
					if err != nil {
						if err.Error() == "Bad Request: message can't be deleted" {
							_, err = msg.ReplyText(err.Error())
							return err
						}
					}
				}
				err = logger.SendLog(b, u, "username", "")
				return err
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func picture(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	db := sql.GetPicture(chat.Id)

	if db.Option != "true" {
		return gotgbot.EndGroups{}
	}
	if chat_status.IsUserAdmin(chat, msg.From.Id) == true {
		return gotgbot.EndGroups{}
	}

	photo, err := user.GetProfilePhotos(0, 0)
	err_handler.HandleErr(err)

	if photo != nil && photo.TotalCount == 0 {
		bantime := extraction.ExtractTime(b, msg, sql.GetSetting(chat.Id).Time)
		replytext := GetStringf(msg.Chat.Id, "modules/listener.go:173",
			map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": db.Action,
				"4": strconv.Itoa(user.Id)})

		kb := make([][]ext.InlineKeyboardButton, 1)
		kb[0] = make([]ext.InlineKeyboardButton, 1)
		kb[0][0] = ext.InlineKeyboardButton{Text: GetString(chat.Id, "modules/listener.go:179"),
			CallbackData: fmt.Sprintf("pmute_%v", user.Id)}

		kbk := make([][]ext.InlineKeyboardButton, 1)
		kbk[0] = make([]ext.InlineKeyboardButton, 1)
		kbk[0][0] = ext.InlineKeyboardButton{Text: GetString(chat.Id, "modules/listener.go:184"),
			CallbackData: fmt.Sprintf("pban_%v", user.Id)}

		// Implementing The Chosen Action To The Target
		if db.Action == "mute" {
			restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = bantime
			_, err = restrictSend.Send()
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					_, err = b.SendMessage(chat.Id, err.Error())
					return err
				}
			}
			// Send Message
			reply := b.NewSendableMessage(chat.Id, replytext)
			reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
			reply.ParseMode = parsemode.Html
			reply.ReplyToMessageId = msg.MessageId
			_, err = reply.Send()
			if err != nil {
				if err.Error() == "Bad Request: reply message not found" {
					reply := b.NewSendableMessage(chat.Id, replytext)
					reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
					reply.ParseMode = parsemode.Html
					_, err = reply.Send()
					return err
				} else {
					err_handler.HandleErr(err)
				}
			}
		} else if db.Action == "kick" {
			restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
			_, err = restrictSend.Send()
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					_, err = b.SendMessage(chat.Id, err.Error())
					return err
				}
			}
			_, err = b.UnbanChatMember(chat.Id, user.Id)
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					_, err = b.SendMessage(chat.Id, err.Error())
					return err
				}
			}
			// Send Message
			reply := b.NewSendableMessage(chat.Id, replytext)
			reply.ParseMode = parsemode.Html
			reply.ReplyToMessageId = msg.MessageId
			_, err = reply.Send()
			return err
		} else if db.Action == "ban" {
			restrictSend := b.NewSendableKickChatMember(chat.Id, user.Id)
			restrictSend.UntilDate = -1
			_, err = restrictSend.Send()
			if err != nil {
				if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
					_, err = b.SendMessage(chat.Id, err.Error())
					return err
				}
			}
			// Send Message
			reply := b.NewSendableMessage(chat.Id, replytext)
			reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kbk}
			reply.ParseMode = parsemode.Html
			reply.ReplyToMessageId = msg.MessageId
			_, err = reply.Send()
			if err != nil {
				if err.Error() == "Bad Request: reply message not found" {
					reply := b.NewSendableMessage(chat.Id, replytext)
					reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kbk}
					reply.ParseMode = parsemode.Html
					_, err = reply.Send()
					return err
				} else {
					err_handler.HandleErr(err)
				}
			}
		}
		// Delete His/Her Message(s)
		if db.Deletion == "true" {
			_, err = msg.Delete()
			if err != nil {
				if err.Error() == "Bad Request: message can't be deleted" {
					_, err = msg.ReplyText(err.Error())
					return err
				}
			}
		}
		// Send Log
		err = logger.SendLog(b, u, "picture", "")
		return err
	}
	return gotgbot.ContinueGroups{}
}

func verify(b ext.Bot, u *gotgbot.Update) error {
	var err error
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	db := sql.GetVerify(chat.Id)

	if db.Option != "true" {
		return gotgbot.EndGroups{}
	}
	if chat_status.IsUserAdmin(chat, msg.From.Id) == true {
		return gotgbot.EndGroups{}
	}

	if msg != nil {
		bantime := extraction.ExtractTime(b, msg, sql.GetSetting(chat.Id).Time)
		replytext := GetStringf(msg.Chat.Id, "modules/listener.go:298",
			map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": chat.Title, "4": strconv.Itoa(user.Id)})

		kb := make([][]ext.InlineKeyboardButton, 1)
		kb[0] = make([]ext.InlineKeyboardButton, 1)
		kb[0][0] = ext.InlineKeyboardButton{Text: GetString(chat.Id, "modules/listener.go:303"),
			CallbackData: fmt.Sprintf("wlcm_%v", user.Id)}

		restrictSend := b.NewSendableRestrictChatMember(chat.Id, user.Id)
		restrictSend.UntilDate = bantime
		_, err = restrictSend.Send()
		if err != nil {
			if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
				_, err = b.SendMessage(chat.Id, err.Error())
				return err
			}
		}

		reply := b.NewSendableMessage(chat.Id, replytext)
		reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
		reply.ParseMode = parsemode.Html
		reply.ReplyToMessageId = msg.MessageId
		_, err = reply.Send()
		if err != nil {
			if err.Error() == "Bad Request: reply message not found" {
				reply := b.NewSendableMessage(chat.Id, replytext)
				reply.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
				reply.ParseMode = parsemode.Html
				_, err = reply.Send()
				return err
			}
		}
		// Delete His/Her Message(s)
		if db.Deletion == "true" {
			_, err = msg.Delete()
			if err != nil {
				if err.Error() == "Bad Request: message can't be deleted" {
					_, err = msg.ReplyText(err.Error())
					return err
				}
			}
		}
		err = logger.SendLog(b, u, "welcome", "")
		return err
	}
	return gotgbot.ContinueGroups{}
}

func spam(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat_status.IsUserAdmin(chat, msg.From.Id) == true {
		return gotgbot.EndGroups{}
	}
	if chat_status.IsBotAdmin(chat, nil) == false {
		return gotgbot.EndGroups{}
	}

	if msg != nil {
		if chat.Type == "supergroup" {
			if sql.GetEnforceGban(chat.Id).Option == "true" {
				go func() {
					r := &spammer{}
					response, err := http.Get(fmt.Sprintf("https://combot.org/api/cas/check?user_id=%v", user.Id))
					err_handler.HandleErr(err)
					body, err := ioutil.ReadAll(response.Body)
					err_handler.HandleErr(err)
					err = json.Unmarshal(body, &r)
					err_handler.HandleErr(err)
					if r.Status == true {
						err = spamfunc(b, u)
						err_handler.HandleErr(err)
						err = logger.SendLog(b, u, "spam", "CAS Banned (Powered By CAS)")
						err_handler.HandleErr(err)
					}
				}()

				ban := sql.GetUserSpam(user.Id)
				if ban != nil {
					err := spamfunc(b, u)
					err_handler.HandleErr(err)
					err = logger.SendLog(b, u, "spam", ban.Reason)
					return err
				}
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func update(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if msg != nil {
		go sql.UpdateUser(user.Id, user.Username, user.FirstName)
		go sql.UpdateChat(strconv.Itoa(chat.Id), chat.Title, chat.Type, chat.InviteLink)

		if msg.ForwardFrom != nil {
			usr := msg.ForwardFrom
			go sql.UpdateUser(usr.Id, usr.Username, usr.FirstName)
		}

		if sql.GetVerify(chat.Id) == nil {
			go sql.UpdateVerify(chat.Id, "true", "-", "true")
		}
		if sql.GetUsername(chat.Id) == nil {
			go sql.UpdateUsername(chat.Id, "true", "mute", "-", "true")
		}
		if sql.GetPicture(chat.Id) == nil {
			go sql.UpdatePicture(chat.Id, "true", "mute", "-", "true")
		}
		if sql.GetSetting(chat.Id) == nil {
			go sql.UpdateSetting(chat.Id, "5m", "true")
		}
		if sql.GetEnforceGban(chat.Id) == nil {
			go sql.UpdateEnforceGban(chat.Id, "true")
		}
		if sql.GetLang(chat.Id) == nil {
			sql.REDIS.Set(fmt.Sprintf("lang_%v", chat.Id), "ID", 0)
			sql.REDIS.BgSave()
			go sql.UpdateLang(chat.Id, "ID")
		}
	}
	return gotgbot.ContinueGroups{}
}

func usernamequery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			data, _ := regexp.MatchString("^umute_\\d+$", msg.Data)
			data2, _ := regexp.MatchString("^uba_\\d+$", msg.Data)
			if data == true {
				if strings.Split(msg.Data, "umute_")[1] == strconv.Itoa(user.Id) {
					if user.Username != "" {
						_, err = b.AnswerCallbackQueryText(msg.Id, GetString(chat.Id, "modules/listener.go:441"), true)
						if err != nil {
							_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
							return err
						}
						_, err = msg.Message.Delete()
						if err != nil {
							_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
							return err
						}
						_, err = b.UnRestrictChatMember(chat.Id, user.Id)
						return err
					} else {
						_, err = b.AnswerCallbackQueryText(msg.Id, GetString(chat.Id, "modules/listener.go:454"), true)
						return err
					}
				} else {
					_, err = b.AnswerCallbackQueryText(msg.Id, GetString(chat.Id, "modules/listener.go:458"), true)
					return err
				}
			} else if data2 == true {
				if chat_status.IsUserAdmin(chat, user.Id) == true {
					i, _ := strconv.Atoi(strings.Split(msg.Data, "uba_")[1])
					_, err = b.UnbanChatMember(chat.Id, i)
					_, err = b.AnswerCallbackQueryText(msg.Id, "Unbanned!", true)
					if err != nil {
						_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
						return err
					}
					_, err = msg.Message.Delete()
					if err != nil {
						_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
						return err
					}

					return err
				}
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func picturequery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			photo, _ := user.GetProfilePhotos(0, 0)
			data, _ := regexp.MatchString("^pmute_\\d+$", msg.Data)
			data2, _ := regexp.MatchString("^pban_\\d+$", msg.Data)
			if data == true {
				if strings.Split(msg.Data, "pmute_")[1] == strconv.Itoa(user.Id) {
					if photo.TotalCount != 0 {
						_, err = b.AnswerCallbackQueryText(msg.Id, GetString(chat.Id, "modules/listener.go:498"), true)
						if err != nil {
							_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
							return err
						}
						_, err = msg.Message.Delete()
						if err != nil {
							_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
							return err
						}
						_, err = b.UnRestrictChatMember(chat.Id, user.Id)
						return err
					} else {
						_, err = b.AnswerCallbackQueryText(msg.Id, GetString(chat.Id, "modules/listener.go:511"), true)
						return err
					}
				} else {
					_, err = b.AnswerCallbackQueryText(msg.Id, GetString(chat.Id, "modules/listener.go:515"), true)
					return err
				}
			} else if data2 == true {
				if chat_status.IsUserAdmin(chat, user.Id) == true {
					i, _ := strconv.Atoi(strings.Split(msg.Data, "pban_")[1])
					_, err = b.UnbanChatMember(chat.Id, i)
					_, err = b.AnswerCallbackQueryText(msg.Id, "Unbanned!", true)
					if err != nil {
						_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
						return err
					}
					_, err = msg.Message.Delete()
					if err != nil {
						_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
						return err
					}

					return err
				}
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func verifyquery(b ext.Bot, u *gotgbot.Update) error {
	var err error
	msg := u.CallbackQuery
	user := msg.From
	chat := msg.Message.Chat

	if msg != nil {
		if chat.Type == "supergroup" {
			data, _ := regexp.MatchString("^wlcm_\\d+$", msg.Data)
			if data == true {
				if strings.Split(msg.Data, "wlcm_")[1] == strconv.Itoa(user.Id) {
					_, err = b.AnswerCallbackQueryText(msg.Id, GetString(chat.Id, "modules/listener.go:552"), true)
					if err != nil {
						_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
						return err
					}
					_, err = msg.Message.Delete()
					if err != nil {
						_, err = b.AnswerCallbackQueryText(msg.Id, err.Error(), true)
						return err
					}
					_, err = b.UnRestrictChatMember(chat.Id, user.Id)
					return err

				} else {
					_, err = b.AnswerCallbackQueryText(msg.Id, GetString(chat.Id, "modules/listener.go:566"), true)
					return err
				}
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func spamfunc(b ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	db := sql.GetSetting(chat.Id)
	txtBan := GetStringf(chat.Id, "modules/listener.go:580",
		map[string]string{"1": strconv.Itoa(user.Id), "2": user.FirstName, "3": strconv.Itoa(user.Id)})

	_, err := msg.ReplyHTMLf(txtBan)
	if err != nil {
		if err.Error() == "Bad Request: reply message not found" {
			_, err = b.SendMessageHTML(chat.Id, txtBan)
			return err
		}
	}
	_, err = b.KickChatMember(chat.Id, user.Id)
	if err != nil {
		if err.Error() == "Bad Request: not enough rights to restrict/unrestrict chat member" {
			_, err = b.SendMessage(chat.Id, err.Error())
			return err
		}
	}

	if db.Deletion == "true" {
		_, err = b.DeleteMessage(chat.Id, msg.MessageId)
		if err != nil {
			if err.Error() == "Bad Request: message can't be deleted" {
				_, err = b.SendMessage(chat.Id, err.Error())
				return err
			}
		}
	}
	return nil
}

func LoadListeners(u *gotgbot.Updater) {
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, update))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, spam))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, username))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, picture))
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.NewChatMembers(), verify))
	u.Dispatcher.AddHandler(handlers.NewCallback(regexp.MustCompile("^(umute|uba)_\\d+$").String(), usernamequery))
	u.Dispatcher.AddHandler(handlers.NewCallback(regexp.MustCompile("^(pmute|pban)_\\d+$").String(), picturequery))
	u.Dispatcher.AddHandler(handlers.NewCallback(regexp.MustCompile("^wlcm_\\d+$").String(), verifyquery))
}
