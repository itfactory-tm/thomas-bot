package main

import (
	"log"
	"time"

	"github.com/dghubble/oauth1"

	"github.com/dghubble/go-twitter/twitter"

	"github.com/bwmarrin/discordgo"
)

const whatsupChannel = "697150309482496082"

func postHashtagTweets(s *discordgo.Session) {
	if !c.TwitterEnabled {
		log.Println("Twitter posting disabled")
		return
	}

	config := oauth1.NewConfig(c.TwitterConsumerKey, c.TwitterConsumerSecret)
	token := oauth1.NewToken(c.TwitterAccessToken, c.TwitterAccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	params := &twitter.StreamFilterParams{
		Track:         []string{"#ITFactory"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(params)
	if err != nil {
		log.Println(err)
		return
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if tweet.Retweeted {
			return
		}

		if tweet.User.FollowersCount < 5 {
			// we do not take people with less than 5 followers seriously
			return
		}
		embed := NewEmbed()
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

		_, err := s.ChannelMessageSendEmbed(whatsupChannel, embed.MessageEmbed)
		if err != nil {
			log.Println(err)
		}
	}

	for {
		log.Println("Starting Twitter listener")
		demux.HandleChan(stream.Messages)
		time.Sleep(10 * time.Second) // backoff in case of crash
	}
}
