package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func main() {
	token := os.Getenv("TOKEN")
	b, err := gotgbot.NewBot(token , &gotgbot.BotOpts{
		Client:      http.Client{},
		GetTimeout:  gotgbot.DefaultGetTimeout,
		PostTimeout: gotgbot.DefaultPostTimeout,
	})

	if err != nil {
		panic(err)
	}

	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: nil,
		DispatcherOpts: ext.DispatcherOpts{
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				fmt.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			Panic:       nil,
			ErrorLog:    nil,
			MaxRoutines: 0,
		},
	})
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(
		handlers.NewMessage(
			message.All, handle,
		),
	)

	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	fmt.Printf("%s has been started...\n", b.User.Username)
	updater.Idle()

}

func check_admin(user int64 , b *gotgbot.Bot) bool {
	x , _ := b.GetChatMember(b.User.Id, user)
	if x.GetStatus() == "creator" || x.GetStatus() == "administrator" {
		return true
	}
	return false

}

func regex(text, regex string) bool {
	mizu, _ := regexp.Compile(regex)
	if mizu.MatchString(text) {
		return true
	}
	return false
}

func handle(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message

	if regex(msg.Text, "^/start$") {
		b.SendMessage(ctx.Message.Chat.Id, "Hello i'm Mizuhara", &gotgbot.SendMessageOpts{
			ParseMode: "Markdown",
		})
	}

	if regex(msg.Text, "^/help$") {
		b.SendMessage(ctx.Message.Chat.Id, "Hello i'm Mizuhara", &gotgbot.SendMessageOpts{
			ParseMode: "Markdown",
		})
	}

	if regex(msg.Text, "(?i)mizu") {
		b.SendMessage(ctx.Message.Chat.Id, "Hello User-Kun", &gotgbot.SendMessageOpts{
			ParseMode: "Markdown",
		})
	}

	if regex(msg.Text, "(?i)mizu ban") && check_admin(msg.From.Id, b) {

		_, err := b.BanChatMember(ctx.Message.Chat.Id, msg.ReplyToMessage.From.Id, &gotgbot.BanChatMemberOpts{
			UntilDate: 0,
		})

		if err != nil {
			b.SendMessage(ctx.Message.Chat.Id, string(err.Error()), &gotgbot.SendMessageOpts{
				ParseMode: "Markdown",
			})
		}

		b.SendMessage(ctx.Message.Chat.Id, "User-Kun has been banned", &gotgbot.SendMessageOpts{})

		fmt.Println(msg.Text)

	}

	if regex(msg.Text, "(?i)mizu unban") && check_admin(msg.From.Id, b) {
		user, _ := strconv.ParseInt(strings.Split(msg.Text, " ")[2], 10, 64)
		fmt.Println(user)
		_, err := b.UnbanChatMember(ctx.Message.Chat.Id, user, &gotgbot.UnbanChatMemberOpts{
			OnlyIfBanned: true,
		})

		if err != nil {
			b.SendMessage(ctx.Message.Chat.Id, err.Error(), &gotgbot.SendMessageOpts{
				ParseMode: "Markdown",
			})
		}
	}

	return nil

}
