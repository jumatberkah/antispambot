package chat_status

import (
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/helpers/function"
	"strconv"
)

func CanDelete(chat *ext.Chat, botId int) bool {
	k, err := chat.GetMember(botId)
	err_handler.HandleErr(err)
	return k.CanDeleteMessages
}

func IsOwner(userId int) bool {
	for _, user := range bot.BotConfig.SudoUsers {
		if user == strconv.Itoa(userId) {
			return true
		}
	}
	return false
}

func IsUserAdmin(chat *ext.Chat, userId int, member *ext.ChatMember) bool {
	if chat.Type == "private" || function.Contains(bot.BotConfig.SudoUsers, strconv.Itoa(userId)) {
		return true
	}
	if member == nil {
		mem, err := chat.GetMember(userId)
		err_handler.HandleErr(err)
		member = mem
	}
	if member.Status == "administrator" || member.Status == "creator" {
		return true
	} else {
		return false
	}
}

func IsBotAdmin(chat *ext.Chat, member *ext.ChatMember) bool {
	if chat.Type == "private" {
		return true
	}
	if member == nil {
		mem, err := chat.GetMember(chat.Bot.Id)
		err_handler.HandleErr(err)
		if mem == nil {
			return false
		}
		member = mem

	}
	if member.Status == "administrator" || member.Status == "creator" {
		return true
	} else {
		return false
	}
}

func RequireBotAdmin(chat *ext.Chat, msg *ext.Message) bool {
	if !IsBotAdmin(chat, nil) {
		_, err := msg.ReplyText("Anda harus menjadikan saya administrator untuk melakukannya.")
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func RequireUserAdmin(chat *ext.Chat, msg *ext.Message, userId int, member *ext.ChatMember) bool {
	if !IsUserAdmin(chat, userId, member) {
		_, err := msg.ReplyText("Anda harus menjadi administrator untuk melakukannya.")
		err_handler.HandleErr(err)
		return false
	}
	return true
}

func IsUserInChat(chat *ext.Chat, userId int) bool {
	member, err := chat.GetMember(userId)
	err_handler.HandleErr(err)
	if member.Status == "left" || member.Status == "kicked" {
		return false
	} else {
		return true
	}
}

func CanRestrict(bot ext.Bot, chat *ext.Chat) bool {
	botChatMember, err := chat.GetMember(bot.Id)
	err_handler.HandleErr(err)
	if !botChatMember.CanRestrictMembers {
		_, err := bot.SendMessage(chat.Id, "Saya tidak bisa membatasi seseorang. "+
			"Pastikan saya admin dan memiliki semua izin")
		err_handler.HandleErr(err)
		return false
	}
	return true
}
