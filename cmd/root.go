package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"gorm.io/gorm"
)

const (
	defaultDBName = "emlyon-tgbot"
)

var options struct {
	verbose     bool
	debug       bool
	generateDoc bool
	dbTruncate  bool
	dbName      string
}

var (
	db  *gorm.DB
	log = logrus.New()
)

// RootCmd is the root command for limo
var RootCmd = &cobra.Command{
	Use:   "emlyon-tgbot",
	Short: "A smart telegram bot about EM Lyon news and knowledge base.",
	Long:  `A smart telegram bot about EM Lyon news and knowledge base running with the latest machine learning technologies.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if options.generateDoc {
		err := doc.GenMarkdownTree(RootCmd, "./docs")
		if err != nil {
			log.Fatal(err)
		}
		// todo: generate changelog and todo from code
	}
}

func init() {
	flags := RootCmd.PersistentFlags()
	flags.StringVarP(&options.dbName, "db-name", "n", defaultDBName, "database name.")
	flags.BoolVarP(&options.dbTruncate, "truncate", "t", false, "truncate database.")
	flags.BoolVarP(&options.generateDoc, "generate-doc", "g", false, "generate documentation.")
	flags.BoolVarP(&options.debug, "debug", "d", false, "debug mode.")
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose output.")
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func ensureDir(dirname string) error {
	log.Infof("ensureDir(%s)\n", dirname)
	st, err := os.Stat(dirname)
	if err != nil {
		log.Infof("ensureDir creating %s\n", dirname)
		err := os.MkdirAll(dirname, 0755)
		if err != nil {
			return err
		}
		st, _ = os.Stat(dirname)
	}
	var inode uint64 = 0
	stat, ok := st.Sys().(*syscall.Stat_t)
	if ok {
		inode = stat.Ino
	}
	log.Infof("ensureDir(%s) DONE inode %d\n", dirname, inode)
	return nil
}
