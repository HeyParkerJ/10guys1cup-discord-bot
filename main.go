package main

import (
	"context"
	"os"
	// "regexp"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.InfoLevel,
}

var noCtx = context.Background()

// checkErr logs errors if not nil, along with a user-specified trace
func checkErr(err error, trace string) {
	if err != nil {
		log.WithFields(logrus.Fields{
			"trace": trace,
		}).Error(err)
	}
}

// handleMsg is a basic command handler
func handleMsg(s disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message

	// Declare regex matchers for contain checks
	// r, _ := regexp.Compile("tilt")

	if msg.Content == "ping" {
		_, err := msg.Reply(noCtx, s, "pong")
		checkErr(err, "ping command")
	} else if strings.Contains("tilt", msg.Content) {
		reply := "Dude FUCK THIS YEAR. Honestly garbage. Playing the fucking highest scorer every fucking week but still putting up points while the rest of the shit teams get wins. Fuck this dude. I'll auto draft next year. Full fuckin tilt."
		_, err := msg.Reply(noCtx, s, reply)
		checkErr(err, "tilt response")
	} else {
		return
	}
}

func init() {
	// In non prod, will load the contents of .env to environment variables
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func main() {
	// var envs map[string]string
	// envs, err := godotenv.Read(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	discordToken := os.Getenv("DISCORD_TOKEN")

	const prefix = "!"

	client := disgord.New(disgord.Config{
		ProjectName: "DadBot",
		BotToken:    discordToken,
		Logger:      log,
		RejectEvents: []string{
			// rarely used, and causes unnecessary spam
			disgord.EvtTypingStart,

			// these require special privilege
			// https://discord.com/developers/docs/topics/gateway#privileged-intents
			disgord.EvtPresenceUpdate,
			disgord.EvtGuildMemberAdd,
			disgord.EvtGuildMemberUpdate,
			disgord.EvtGuildMemberRemove,
		},
		// ! Non-functional due to a current bug, will be fixed.
		Presence: &disgord.UpdateStatusPayload{
			Game: &disgord.Activity{
				Name: "write " + prefix + "ping",
			},
		},
		DMIntents: disgord.IntentDirectMessages | disgord.IntentDirectMessageReactions | disgord.IntentDirectMessageTyping,
		// comment out DMIntents if you do not want the bot to handle direct messages

	})

	defer client.Gateway().StayConnectedUntilInterrupted()

	logFilter, _ := std.NewLogFilter(client)
	filter, _ := std.NewMsgFilter(context.Background(), client)
	// Gonna go without a prefix for now
	// filter.SetPrefix(prefix)

	// create a handler and bind it to new message events
	// thing about the middlewares are whitelists or passthrough functions.
	client.Gateway().WithMiddleware(
		filter.NotByBot, // ignore bot messages
		// filter.HasPrefix,   // message must have the given prefix
		logFilter.LogMsg, // log command message
		// filter.StripPrefix, // remove the command prefix from the message
	).MessageCreate(handleMsg)

	// create a handler and bind it to the bot init
	// dummy log print
	client.Gateway().BotReady(func() {
		log.Info("Bot is ready!")
	})
}
