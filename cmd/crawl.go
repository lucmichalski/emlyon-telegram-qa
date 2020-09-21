package cmd

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/lucmichalski/emlyon-telegram-qa/pkg/models"
)

const (
	sitemapPattern = `https://www.em-lyon.com/%s/sitemap/xml`
)

var (
	cacheReset    bool
	cacheDisabled bool
	cacheDir      string

	languages []string

	parallelJobs   int
	validLanguages = []string{"en", "fr", "cn"}
	validDomains   = []string{"www.em-lyon.com", "em-lyon.com"}
)

// CrawlCmd allows to crawl the website
var CrawlCmd = &cobra.Command{
	Use:     "crawl",
	Aliases: []string{"c", "crawler"},
	Short:   "Crawl pages from the website",
	Long:    "Crawl pages from the website in all languages or specific one...",
	Run: func(cmd *cobra.Command, args []string) {

		// open the database
		var err error
		db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", options.dbName)), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		// Migrate the schema
		db.AutoMigrate(&models.Page{})
		db.AutoMigrate(&models.Document{})

		// if cacheReset {
		// 	os.Remove(cacheDir)
		// }

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

		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("visiting", r.URL)
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
	CrawlCmd.Flags().StringVarP(&cacheDir, "cache-dir", "", "./shared/data/emlyon", "Cache output path.")
	CrawlCmd.Flags().BoolVarP(&cacheReset, "cache-reset", "", false, "Reset http cache")
	CrawlCmd.Flags().BoolVarP(&cacheDisabled, "no-cache", "", false, "Disable http cache")
	CrawlCmd.Flags().IntVarP(&parallelJobs, "jobs", "j", 4, "Parallel jobs")
	RootCmd.AddCommand(CrawlCmd)
}