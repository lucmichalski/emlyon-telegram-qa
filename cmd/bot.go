package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/abadojack/whatlanggo"
	"github.com/dghubble/sling"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"

	"github.com/lucmichalski/emlyon-telegram-qa/pkg/models"
)

var (
	qaTopKReader     int
	qaTopKRetriever  int
	qaServerAddress  string
	qaServerEndpoint string

	tgToken   string
	tgTimeout = 60
)

// CrawlCmd allows to crawl the website
var ChatbotCmd = &cobra.Command{
	Use:     "chatbot",
	Aliases: []string{"c", "chatbot"},
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

			pp.Println("LanguageCode", update.Message.From.LanguageCode)
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			info := whatlanggo.Detect(update.Message.Text)
			if options.debug {
				pp.Println("======> Language:", info.Lang.String(), " Script:", whatlanggo.Scripts[info.Script], " Confidence: ", info.Confidence)
			}

			var indexQASuffix, dontKnowMsg string
			switch info.Lang.Iso6391() {

			case "fr":
				indexQASuffix = "fr"
				dontKnowMsg = "Désolé, je n'ai pas de réponse à cette question..."
				qaServerAddress = "http://localhost:8007"

			case "en":
				fallthrough

			default:
				indexQASuffix = "en"
				dontKnowMsg = "Sorry, I do not have an answer to that question..."
				qaServerAddress = "http://localhost:8006"
			}

			params := &models.HaystackParams{
				Index:         fmt.Sprintf("emlyon-%s", indexQASuffix),
				Question:      update.Message.Text,
				TopKReader:    qaTopKReader,
				TopKRetriever: qaTopKRetriever,
			}

			httpClient := &http.Client{}
			req, err := sling.New().Get(qaServerAddress).Path(qaServerEndpoint).QueryStruct(params).Request()
			resp, err := httpClient.Do(req)
			if err != nil {
				log.Warn(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Warn(err)
			}

			var qaResponse models.Haystack
			err = json.Unmarshal(body, &qaResponse)
			if err != nil {
				log.Debugln("question:", update.Message.Text)
			}

			if len(qaResponse.Answers) == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, dontKnowMsg)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			} else {

				canonicalLinks := make(map[string]string, len(qaResponse.Answers))
				for _, answer := range qaResponse.Answers {

					// rawMsg := fmt.Sprintf(`Page:
					// 	%s

					// 	Answer:
					// 	%s

					// 	Context:
					// 	%s

					// 	Summary:
					// 	%s`, answer.Meta.CanonicalLink, answer.Answer, answer.Context, answer.Meta.Summary)
					rawMsg := fmt.Sprintf(`%s
					----------
					%s`, answer.Meta.CanonicalLink, answer.Meta.Summary)
					fmt.Println("rawMsg", rawMsg)
					url := answer.Meta.CanonicalLink
					canonicalLinks[url] = rawMsg
				}

				for _, canonicalLink := range canonicalLinks {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, canonicalLink)
					// msg.ReplyToMessageID = update.Message.MessageID
					// msg.ParseMode = "Markdown"
					bot.Send(msg)
				}
			}

		}

	},
}

func init() {
	// https://t.me/octaveqa_bot
	ChatbotCmd.Flags().StringVarP(&qaServerEndpoint, "qa-endoint", "", "/query", "Hasytack QA endpoint path.")
	ChatbotCmd.Flags().IntVarP(&qaTopKRetriever, "top-k-retriever", "", 20, "Hasytack QA top-k-retriever value")
	ChatbotCmd.Flags().IntVarP(&qaTopKReader, "top-k-reader", "", 2, "Hasytack QA top-k-reader value")
	ChatbotCmd.Flags().StringVarP(&tgToken, "token-api", "", os.Getenv("EMLYON_TELEGRAM_TOKEN"), "Telegram token api")
	ChatbotCmd.Flags().IntVarP(&tgTimeout, "timeout", "", 60, "Reset http cache")
	RootCmd.AddCommand(ChatbotCmd)
}
