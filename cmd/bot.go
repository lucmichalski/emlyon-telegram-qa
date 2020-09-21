package cmd

import (
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/cobra"
)

var (
	tgToken   string
	tgTimeout = 60
)

// CrawlCmd allows to crawl the website
var ChatbotCmd = &cobra.Command{
	Use:     "chatbot",
	Aliases: []string{"c", "crawler"},
	Short:   "Telegram chatbot service",
	Long:    "Telegram chatbot service with Question-Answering system.",
	Run: func(cmd *cobra.Command, args []string) {

		bot, err := tgbotapi.NewBotAPI(tgToken)
		if err != nil {
			log.Panic(err)
		}
		bot.Debug = options.debug

		log.Printf("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = tgTimeout

		updates, err := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

	},
}

func init() {
	ChatbotCmd.Flags().StringVarP(&tgToken, "token-api", "", os.Getenv("TELEGRAM_TOKEN"), "Telegram token api")
	ChatbotCmd.Flags().IntVarP(&tgTimeout, "timeout", "", 60, "Reset http cache")
	ChatbotCmd.Flags().BoolVarP(&cacheDisabled, "no-cache", "", false, "Disable http cache")
	ChatbotCmd.Flags().IntVarP(&parallelJobs, "jobs", "j", 4, "Parallel jobs")
	RootCmd.AddCommand(ChatbotCmd)
}
