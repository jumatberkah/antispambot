package main

import (
	"github.com/PaulSonOfLars/goloc"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/jumatberkah/antispambot/bot/modules"
	"github.com/jumatberkah/antispambot/bot/modules/sql"
	"log"
)

func main() {
	// initiation
	goloc.LoadAll("en-GB")
	goloc.LoadAll("ID")
	updater, err := gotgbot.NewUpdater(bot.BotConfig.ApiKey)
	err_handler.FatalError(err)

	// registering handlers
	modules.LoadLang(updater)
	modules.LoadPm(updater)
	modules.LoadSetting(updater)
	modules.LoadSettingPanel(updater)
	modules.LoadAdmins(updater)
	modules.LoadListeners(updater)

	// start clean polling / webhook
	if bot.BotConfig.WebhookUrl != "" {
		log.Print("Using Webhook...")
		var web gotgbot.Webhook
		web.URL = bot.BotConfig.WebhookUrl
		web.MaxConnections = 40
		web.Serve = "localhost"
		web.ServePort = bot.BotConfig.WebhookPort
		_, err = updater.SetWebhook(bot.BotConfig.WebhookPath, web)
		err_handler.HandleErr(err)
		updater.StartWebhook(web)
	} else {
		log.Print("Using Long Polling...")
		_ = updater.StartPolling()
	}

	// connect to db
	sql.InitDb(*updater.Bot, nil)

	// wait
	updater.Idle()
}
