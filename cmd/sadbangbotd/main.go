package main

import (
	"os"
	"regexp"

	"github.com/ChimeraCoder/anaconda"
	log "github.com/Sirupsen/logrus"
	"github.com/ralfonso/sadbangbot/streaming"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

const (
	appName = "sadbangbotd"
)

var (
	twitterConsumerKeyFlag = cli.StringFlag{
		Name:   "twitter.consumer.key",
		EnvVar: "SADBANGBOT_TWITTER_CONSUMER_KEY",
	}
	twitterConsumerSecretFlag = cli.StringFlag{
		Name:   "twitter.consumer.secret",
		EnvVar: "SADBANGBOT_TWITTER_CONSUMER_SECRET",
	}
	twitterAccessTokenFlag = cli.StringFlag{
		Name:   "twitter.access.token",
		EnvVar: "SADBANGBOT_TWITTER_ACCESS_TOKEN",
	}
	twitterAccessTokenSecretFlag = cli.StringFlag{
		Name:   "twitter.access.token.secret",
		EnvVar: "SADBANGBOT_TWITTER_ACCESS_TOKEN_SECRET",
	}
)

var (
	globalFlags = []cli.Flag{
		twitterConsumerKeyFlag,
		twitterConsumerSecretFlag,
		twitterAccessTokenFlag,
		twitterAccessTokenSecretFlag,
	}

	sadRe = regexp.MustCompile("[[:punct:]][ ]{0,2}Sad!$")
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Action = mainAction
	app.Flags = globalFlags

	app.Run(os.Args)
}

func mainAction(c *cli.Context) error {
	log.SetLevel(log.DebugLevel)
	api := twitterClient(c)
	ctx := context.Background()

	cCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// enable the anaconda logger so we get some debugging info
	api.SetLogger(anaconda.BasicLogger)

	server := streaming.NewStreamingServer(cCtx, api, filterFunc)
	// server.Publisher = func(status string, _ url.Values) (anaconda.Tweet, error) {
	// 	log.WithFields(log.Fields{
	// 		"tweet.text": status,
	// 	}).Info("would have tweeted!")
	// 	return anaconda.Tweet{Text: status}, nil
	// }
	log.Info("starting streaming server")
	server.Start()
	return nil
}

func filterFunc(tweet anaconda.Tweet) bool {
	// log.WithFields(log.Fields{
	// 	"tweet.text": tweet.Text,
	// }).Debug("checking tweet")
	return sadRe.Match([]byte(tweet.Text))
}

func mustGlobalString(c *cli.Context, flag cli.StringFlag) string {
	s := c.GlobalString(flag.Name)
	if s == "" {
		log.Fatalf("you must set a value for %s", flag.Name)
	}
	return s
}

func twitterClient(c *cli.Context) *anaconda.TwitterApi {
	consumerKey := mustGlobalString(c, twitterConsumerKeyFlag)
	consumerSecret := mustGlobalString(c, twitterConsumerSecretFlag)
	accessToken := mustGlobalString(c, twitterAccessTokenFlag)
	accessTokenSecret := mustGlobalString(c, twitterAccessTokenSecretFlag)

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	return anaconda.NewTwitterApi(accessToken, accessTokenSecret)
}
