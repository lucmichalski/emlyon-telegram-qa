package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/advancedlogic/GoOse"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"

	"github.com/lucmichalski/emlyon-telegram-qa/pkg/articletext"
	"github.com/lucmichalski/emlyon-telegram-qa/pkg/models"
)

const (
	sitemapPattern = `https://www.em-lyon.com/%s/sitemap/xml`
)

var (
	cacheReset    bool
	cacheDisabled bool
	cacheDir      string

	dataReset bool
	dataDir   string

	languages []string

	parallelJobs   int
	validLanguages = []string{"en", "fr"}
	validDomains   = []string{"www.em-lyon.com", "em-lyon.com", "learninghub.em-lyon.com"}
)

// CrawlCmd allows to crawl the website
var CrawlCmd = &cobra.Command{
	Use:     "crawl",
	Aliases: []string{"c", "crawler"},
	Short:   "Crawl pages from the website",
	Long:    "Crawl pages from the website in all languages or specific one...",
	Run: func(cmd *cobra.Command, args []string) {

		// open the database
		// var err error
		// db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", options.dbName)), &gorm.Config{})
		// if err != nil {
		// 	panic("failed to connect database")
		// }

		// // Migrate the schema
		// db.AutoMigrate(&models.Page{})
		// db.AutoMigrate(&models.Document{})

		// if cacheReset {
		// 	os.Remove(cacheDir)
		// }

		if dataReset {
			dir, err := ioutil.ReadDir(dataDir)
			if err != nil {
				log.Fatal(err)
			}
			for _, d := range dir {
				if err := os.RemoveAll(path.Join([]string{dataDir, d.Name()}...)); err != nil {
					log.Fatal(err)
				}
			}
		}

		// Set to empty the cache dir to disable it
		if cacheDisabled {
			cacheDir = ""
		}

		// Create a Collector specifically for Shopify
		c := colly.NewCollector(
			colly.AllowedDomains(validDomains...),
			colly.CacheDir(cacheDir),
		)

		// create a request queue with 2 consumer threads
		q, _ := queue.New(
			parallelJobs, // Number of consumer threads
			&queue.InMemoryQueueStorage{MaxSize: 100000}, // Use default queue storage
		)

		// Create a callback on the XPath query searching for the URLs
		c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
			log.Infof("Enqueuing url '%s'", e.Text)
			q.AddURL(e.Text)
		})

		// Set error handler
		c.OnError(func(r *colly.Response, err error) {
			fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		})

		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			q.AddURL(e.Attr("href"))
		})

		c.OnHTML("html", func(e *colly.HTMLElement) {
			rawHTML, err := e.DOM.Html()
			if err != nil {
				log.Fatalf("error while getting DOM for url='%s'\n", e.Request.Ctx.Get("url"))
			}

			g := goose.New()
			article, _ := g.ExtractFromRawHTML(rawHTML, e.Request.Ctx.Get("url"))
			if options.debug {
				pp.Println("MetaKeywords", article.MetaKeywords)
				pp.Println("MetaLang", article.MetaLang)
				pp.Println("CanonicalLink", article.CanonicalLink)
				pp.Println("CleanedText", article.CleanedText)
				pp.Println("FinalURL", article.FinalURL)
				pp.Println("TopImage", article.TopImage)
				pp.Println("PublishDate", article.PublishDate)
			}
			article.CleanedText = strings.Replace(article.CleanedText, "\n", " ", -1)

			// nb. use article Text to extract page content as goose cleaner is working too efficiently sometimes ^^/
			articleText, err := articletext.GetArticleText(strings.NewReader(rawHTML))
			if err != nil {
				log.Fatalf("error while getting article text for url='%s'\n", e.Request.Ctx.Get("url"))
			}
			articleText = strings.Replace(articleText, "\n", " ", -1)

			if options.debug {
				pp.Println("Text", articleText)
			}

			info := whatlanggo.Detect(articleText)
			if options.debug {
				pp.Println("======> Language:", info.Lang.String(), " Script:", whatlanggo.Scripts[info.Script], " Confidence: ", info.Confidence)
			}

			linkHash := getMD5Hash(article.FinalURL)

			entry := &models.Page{
				LinkHash:           linkHash,
				Title:              article.Title,
				Body:               articleText,
				MetaDescription:    article.MetaDescription,
				MetaKeywords:       article.MetaKeywords,
				MetaLang:           article.MetaLang,
				CanonicalLink:      article.CanonicalLink,
				CleanedText:        article.CleanedText,
				FinalURL:           article.FinalURL,
				TopImage:           article.TopImage,
				Language:           info.Lang.String(),
				LanguageConfidence: info.Confidence,
			}

			if article.PublishDate != nil {
				entry.PublishDate = *article.PublishDate
			}

			jsonStr, err := json.Marshal(&entry)
			if err != nil {
				log.Fatalln("error while marshalling, msg=", err)
			}

			if article.MetaLang == "" {
				article.MetaLang = "unkown"
			}

			prefixDir := fmt.Sprintf("%s/%s", dataDir, article.MetaLang)
			if err := ensureDir(prefixDir); err != nil {
				log.Fatal(err)
			}

			jsonFile := fmt.Sprintf("%s/%s.json", prefixDir, linkHash)
			if options.debug {
				pp.Println("jsonFile:", jsonFile)
			}

			if err := ioutil.WriteFile(jsonFile, jsonStr, os.ModePerm); err != nil {
				log.Fatal(err)
			}

		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("visiting", r.URL)
			r.Ctx.Put("url", r.URL.String())
		})

		// enqueue sitemap url URL
		for _, language := range languages {
			q.AddURL(fmt.Sprintf(sitemapPattern, language))
		}

		// Consume URLs
		q.Run(c)

	},
}

func init() {
	CrawlCmd.Flags().StringSliceVarP(&languages, "languages", "l", []string{"en"}, fmt.Sprintf("Language to crawl (valid: %s)", strings.Join(validLanguages, ",")))
	CrawlCmd.Flags().BoolVarP(&dataReset, "data-reset", "", false, "Reset collected data (json files)")
	CrawlCmd.Flags().StringVarP(&dataDir, "data-dir", "", "./shared/data/emlyon", "Data directory.")
	CrawlCmd.Flags().StringVarP(&cacheDir, "cache-dir", "", "./shared/cache/emlyon", "Cache output path.")
	CrawlCmd.Flags().BoolVarP(&cacheReset, "cache-reset", "", false, "Reset http cache")
	CrawlCmd.Flags().BoolVarP(&cacheDisabled, "no-cache", "", false, "Disable http cache")
	CrawlCmd.Flags().IntVarP(&parallelJobs, "jobs", "j", 4, "Parallel jobs")
	RootCmd.AddCommand(CrawlCmd)
}
