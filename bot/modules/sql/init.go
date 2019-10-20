package sql

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jumatberkah/antispambot/bot"
	"github.com/jumatberkah/antispambot/bot/helpers/err_handler"
	"github.com/lib/pq"
	"log"
)

var SESSION *gorm.DB
var REDIS *redis.Client

func InitDb(b ext.Bot, u *gotgbot.Update) {
	// Parse the URL
	conn, err := pq.ParseURL(bot.BotConfig.SqlUri)
	err_handler.HandleTgErr(b, u, err)

	// Open a session
	db, err := gorm.Open("postgres", conn)
	err_handler.HandleTgErr(b, u, err)
	SESSION = db
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(100)

	// Auto migrate tables
	db.AutoMigrate(&User{}, &Chat{}, &UserSpam{}, &ChatSpam{}, &Setting{}, &Verify{}, &Picture{}, &Username{},
		&EnforceGban{}, &Lang{})
	log.Println("Database has been connected & Auto-migrated database schema")

	// redis
	client := redis.NewClient(&redis.Options{
		Addr:     bot.BotConfig.RedisAddress,
		Password: bot.BotConfig.RedisPassword, // no password set
		DB:       0,                           // use default DB
	})

	REDIS = client
	err = client.Ping().Err()
	if err == nil {
		log.Println("Redis Has Been Connected")
	}
	defer REDIS.BgSave()
	// Output: PONG <nil>
}
