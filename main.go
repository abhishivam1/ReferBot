// Made by @reeshuxd
package main

import (
	"fmt"
	"log"
	"net/http"
	"newbi/db"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

var token = `6479628391:AAG5xPUIz7Y1O4s6LhXgjYlv5hWqiLhj4L4` // Your Bot token here i.e. by @botfather....

func main() {

	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  gotgbot.DefaultAPIURL,
			},
		},
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
	})

	dispatcher := updater.Dispatcher
	dispatcher.AddHandler(handlers.NewCommand("start", Start))
	dispatcher.AddHandler(handlers.NewCommand("join", Join))
	dispatcher.AddHandler(handlers.NewCommand("competition", Competiton))
	dispatcher.AddHandler(handlers.NewCommand("referral", Referral))

	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}

	fmt.Printf("%s has been started...\n", b.User.Username)
	updater.Idle()
}

func Start(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage

	fsub := Fsub(ctx.EffectiveUser.Id, bot)
	if !fsub {
		message.Reply(bot, "Join my channel in order to use me\nLink here...", nil)
		return nil
	}
	if len(ctx.Args()) != 1 {
		// if db.CheckUser(ctx.EffectiveSender.Id()) {
		// 	message.Reply(bot, "You have already starteed the bot!", nil)
		// 	return nil
		// }
		db.Refer_Update(ctx.Args()[1], "e") // Updating user refers....
		message.Reply(bot, fmt.Sprintf("YOu are invited by %s", ctx.Args()[1]), nil)
		db.AddUser(ctx.EffectiveSender.Id())
	} else {
		message.Reply(bot, "Bot is alive", nil)
		db.AddUser(ctx.EffectiveSender.Id())
	}
	return nil
}

func Join(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	fsub := Fsub(ctx.EffectiveUser.Id, bot)
	if !fsub {
		message.Reply(bot, "Join my channel in order to use me\nLink here...", nil)
		return nil
	}
	if ctx.EffectiveChat.Type != "supergroup" && ctx.EffectiveChat.Type != "group" {
		message.Reply(bot, "This bot can only be used in groups", nil)
		return nil
	}

	text := `
Hello %s, Here is your referral link:
Link - https://t.me/%s?start=%s
	`
	message.Reply(bot, fmt.Sprintf(text, ctx.EffectiveSender.FirstName(), bot.Username, ctx.EffectiveSender.Username()), nil)
	return nil
}

func Competiton(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	fsub := Fsub(ctx.EffectiveUser.Id, bot)
	if !fsub {
		message.Reply(bot, "Join my channel in order to use me\nLink here...", nil)
		return nil
	}
	if ctx.EffectiveChat.Type != "supergroup" && ctx.EffectiveChat.Type != "group" {
		message.Reply(bot, "This bot can only be used in groups", nil)
		return nil
	}
	ctx.EffectiveMessage.Reply(bot, "Cashprize: \n🥇: 1200$ \n🥈: 750$ \n🥉: 400$", nil)
	return nil
}

func Referral(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	fsub := Fsub(ctx.EffectiveUser.Id, bot)
	if !fsub {
		message.Reply(bot, "Join my channel in order to use me\nLink here...", nil)
		return nil
	}
	if ctx.EffectiveChat.Type != "supergroup" && ctx.EffectiveChat.Type != "group" {
		message.Reply(bot, "This bot can only be used in groups", nil)
		return nil
	}

	users, err := db.GetUsersByRefersAscending()
	if err != nil {
		message.Reply(bot, err.Error(), nil)
		return nil
	}

	text := "Top Refers:\n"
	for _, user := range users {
		text += fmt.Sprintf(`@%s - %d refers`, user.UserID, user.Refers)
		text += "\n"
	}

	message.Reply(bot, text, &gotgbot.SendMessageOpts{ParseMode: "html"})
	return nil
}

type User struct {
	UserID int64 `bson:"user_id"`
	Refers int64 `bson:"refers"`
}

var FsubChats = []int64{-1001306365800}

func Fsub(user_id int64, bot *gotgbot.Bot) bool {
	chats := FsubChats
	Result := false
	for _, chat := range chats {
		member, err := bot.GetChatMember(chat, user_id, &gotgbot.GetChatMemberOpts{})
		if err != nil {
			bot.SendMessage(
				FsubChats[0],
				fmt.Sprintf("Bot is not admin in the channel - %v", chat),
				&gotgbot.SendMessageOpts{},
			)
			Result = false
		}
		status := member.GetStatus()
		if status == "creator" || status == "member" || status == "administrator" {
			Result = true
			break
		}
		Result = false
	}
	return Result
}

// Made by @reeshuxd
