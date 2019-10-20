package modules

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/jumatberkah/antispambot/bot/helpers/chat_status"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"regexp"
	"strings"
)

func setusername(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat.Type == "supergroup" {
		if chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					go sql.UpdateUsername(chat.Id, "true", "mute", "-", "true")
					_, err := msg.ReplyHTML("<b>Pengaturan telah diperbarui</b>")
					return err
				} else if strings.ToLower(args[0]) == "false" {
					go sql.UpdateUsername(chat.Id, "false", "mute", "-", "true")
					_, err := msg.ReplyHTML("<b>Pengaturan telah diperbarui</b>")
					return err
				} else {
					_, err := msg.ReplyHTML("<b>Masukkan opsi diantara true/false!</b>")
					return err
				}
			} else {
				_, err := msg.ReplyHTML("<b>Masukkan opsi diantara true/false!</b>")
				return err
			}
		}
	}
	return nil
}

func setverify(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat.Type == "supergroup" {
		if chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					go sql.UpdateVerify(chat.Id, "true", "-", "true")
					_, err := msg.ReplyHTML("<b>Pengaturan telah diperbarui</b>")
					return err
				} else if strings.ToLower(args[0]) == "false" {
					go sql.UpdateVerify(chat.Id, "false", "-", "true")
					_, err := msg.ReplyHTML("<b>Pengaturan telah diperbarui</b>")
					return err
				} else {
					_, err := msg.ReplyHTML("<b>Masukkan opsi diantara true/false!</b>")
					return err
				}
			} else {
				_, err := msg.ReplyHTML("<b>Masukkan opsi diantara true/false!</b>")
				return err
			}
		}
	}
	return nil
}

func setenforce(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat.Type == "supergroup" {
		if chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					go sql.UpdateEnforceGban(chat.Id, "true")
					_, err := msg.ReplyHTML("<b>Pengaturan telah diperbarui</b>")
					return err
				} else if strings.ToLower(args[0]) == "false" {
					go sql.UpdateEnforceGban(chat.Id, "false")
					_, err := msg.ReplyHTML("<b>Pengaturan telah diperbarui</b>")
					return err
				} else {
					_, err := msg.ReplyHTML("<b>Masukkan opsi diantara true/false!</b>")
					return err
				}
			} else {
				_, err := msg.ReplyHTML("<b>Masukkan opsi diantara true/false!</b>")
				return err
			}
		}
	}
	return nil
}

func setpicture(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat.Type == "supergroup" {
		if chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
			if len(args) != 0 {
				if strings.ToLower(args[0]) == "true" {
					go sql.UpdatePicture(chat.Id, "true", "mute", "-", "true")
					_, err := msg.ReplyHTML("<b>Pengaturan telah diperbarui</b>")
					return err
				} else if strings.ToLower(args[0]) == "false" {
					go sql.UpdatePicture(chat.Id, "false", "mute", "-", "true")
					_, err := msg.ReplyHTML("<b>Pengaturan telah diperbarui</b>")
					return err
				} else {
					_, err := msg.ReplyHTML("<b>Masukkan opsi diantara true/false!</b>")
					return err
				}
			} else {
				_, err := msg.ReplyHTML("<b>Masukkan opsi diantara true/false!</b>")
				return err
			}
		}
	}
	return nil
}

func settime(b ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	if chat.Type == "supergroup" {
		if chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
			if len(args) != 0 {
				match, err := regexp.MatchString("^\\d+[mhd]", strings.ToLower(args[0]))
				if match == true {
					go sql.UpdateSetting(chat.Id, args[0], "true")
					_, err := msg.ReplyHTML("<b>Pengaturan telah diperbarui</b>")
					return err
				} else {
					_, err = msg.ReplyHTML("<b>Masukkan nilai waktu!</b>")
					return err
				}
			} else {
				_, err := msg.ReplyHTML("<b>Masukkan nilai waktu!</b>")
				return err
			}
		}
	}
	return gotgbot.ContinueGroups{}
}

func LoadSetting(u *gotgbot.Updater) {
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("username", setusername))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("verify", setverify))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("profilepicture", setpicture))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("time", settime))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("enforce", setenforce))
}
