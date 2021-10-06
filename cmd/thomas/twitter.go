package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"github.com/kelseyhightower/envconfig"

	"github.com/spf13/cobra"
)

// TODO: automate these
const whatsupChannel = "697150309482496082"

func init() {
	rootCmd.AddCommand(NewTwitterCmd())
}

type twitterCmdOptions struct {
	Token string

	TwitterConsumerKey       string `envconfig:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret    string `envconfig:"TWITTER_CONSUMER_SECRET"`
	TwitterAccessToken       string `envconfig:"TWITTER_ACCESS_TOKEN"`
	TwitterAccessTokenSecret string `envconfig:"TWITTER_ACCESS_TOKEN_SECRET"`
}

// NewTwitterCmd generates the `twitter` command
func NewTwitterCmd() *cobra.Command {
	s := twitterCmdOptions{}
	c := &cobra.Command{
		Use:     "twitter",
		Short:   "Run the voice server",
		Long:    `This is a separate instance for Twitter posting. This can only run in a single replica`,
		RunE:    s.RunE,
		PreRunE: s.Validate,
	}

	// TODO: switch to viper
	err := envconfig.Process("thomasbot", &s)
	if err != nil {
		log.Fatalf("Error processing envvars: %q\n", err)
	}

	return c
}

func (t *twitterCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if t.Token == "" {
		return errors.New("No token specified")
	}

	if t.TwitterConsumerKey == "" {
		return errors.New("No TWITTER_CONSUMER_KEY specified")
	}
	if t.TwitterConsumerSecret == "" {
		return errors.New("No TWITTER_CONSUMER_SECRET specified")
	}
	if t.TwitterAccessToken == "" {
		return errors.New("No TWITTER_ACCESS_TOKEN specified")
	}
	if t.TwitterAccessTokenSecret == "" {
		return errors.New("No TWITTER_ACCESS_TOKEN_SECRET specified")
	}

	return nil
}

func (t *twitterCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	dg, err := discordgo.New("Bot " + t.Token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	config := oauth1.NewConfig(t.TwitterConsumerKey, t.TwitterConsumerSecret)
	token := oauth1.NewToken(t.TwitterAccessToken, t.TwitterAccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	params := &twitter.StreamFilterParams{
		Track:         []string{"#ITFactory", "#itfactory", "#ITfactory"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(params)
	if err != nil {
		return err
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if tweet.Retweeted {
			return
		}
		if strings.Index(tweet.Text, "RT") == 0 {
			// is retweet
			return
		}

		// trending tweet safeguard
		if tweet.QuotedStatus != nil {
			// quoted tweet, ignoring
			return
		}

		// it is project season relaxing our requirements
		// make sure I do not regret this!

		/*if tweet.User.FollowersCount < 10 {
			// we do not take people with less than 5 followers seriously
			return
		}*/

		if tweet.User.ScreenName == "paardemeis" || tweet.User.ScreenName == "Koosiebz" || tweet.User.ScreenName == "marjolijnWI" || tweet.User.ScreenName == "johnfilibuster" {
			// because lack of a better system
			return
		}

		/*if t, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tweet.User.CreatedAt); err == nil {
			if time.Since(t) < 30*24*time.Hour {
				// accounts need to be 30 days old
				return
			}
		}*/

		_, err := dg.ChannelMessageSend(whatsupChannel, "https://twitter.com/"+tweet.User.ScreenName+"/status/"+tweet.IDStr)
		if err != nil {
			log.Println(err)
		}
		/*embed := embed.NewEmbed()
		embed.AddField("Tweet", tweet.Text)

		images := []string{}

		if tweet.Entities != nil {
			for _, media := range tweet.Entities.Media {
				if media.Type == "photo" {
					images = append(images, media.MediaURLHttps)
				}
			}
		}
		if tweet.ExtendedEntities != nil {
			for _, media := range tweet.ExtendedEntities.Media {
				if media.Type == "photo" || media.Type == "animated_gif" {
					images = append(images, media.MediaURLHttps)
				}
			}
		}
		if len(images) > 0 {
			embed.SetImage(images[0]) // we can only set 1
		}
		embed.SetTitle("@" + tweet.User.ScreenName + ": " + tweet.User.Name)
		embed.SetURL("https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr)
		embed.SetThumbnail(tweet.User.ProfileImageURLHttps)

		_, err := dg.ChannelMessageSendEmbed(whatsupChannel, embed.MessageEmbed)
		if err != nil {
			log.Println(err)
		}*/
	}

	log.Println("Starting Twitter listener")
	demux.HandleChan(stream.Messages)
	stream.Stop()

	return nil
}
