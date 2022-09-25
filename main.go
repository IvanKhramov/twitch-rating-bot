package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDatabase struct {
	WhiteUserName string `yaml:"whiteUserName" env:"WHITE_USER_NAME" env-default:"thekapustaa"`
	ChannelName   string `yaml:"channelName" env:"CHANNEL_NAME" env-default:"hellyeahplay"`
	OAuthToken    string `yaml:"oAuthToken" env:"O_AUTH_TOKEN"`
}

var cfg ConfigDatabase

type topVote struct {
	value int
	count int
}

var (
	usersVotes   = make(map[string]int)
	votesCount   = make(map[int]int)
	ratingStatus = 0
)

func checkAccess(msg twitch.PrivateMessage) bool {
	return msg.User.Name == cfg.WhiteUserName
}

func checkCommand(msg twitch.PrivateMessage, client *twitch.Client) {
	switch strings.ToLower(msg.Message) {
	case "!startrating":
		startRating(client)
	case "!endrating":
		endRating(client)
	}

}

func startRating(client *twitch.Client) {
	client.Say(cfg.ChannelName, "НАЧИНАЕТСЯ ОЦЕНКА АНИМЕ! ПИШИТЕ ЧИСЛО ОТ 1 ДО 10")
	ratingStatus = 1
}

func endRating(client *twitch.Client) {
	var totalSum int
	var top topVote
	for _, value := range usersVotes {
		totalSum += value
		votesCount[value] += 1
	}
	for vote, count := range votesCount {
		if count > top.count {
			top.value, top.count = vote, count
		}
	}

	avgRating := float64(totalSum) / float64(len(usersVotes))

	client.Say(cfg.ChannelName, fmt.Sprintf("Оценка окончена! Средняя оценка - %.1f. Самая частая оценка - %d(проголосовали %d раз).", avgRating, top.value, top.count))
	ratingStatus = 0
	usersVotes = make(map[string]int)
	votesCount = make(map[int]int)
}

func checkRatingMsg(message twitch.PrivateMessage) {
	msg, err := strconv.Atoi(message.Message)
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		return
	}
	if msg > 0 && msg < 11 {
		usersVotes[message.User.Name] = msg
	}
}

func main() {
	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		panic(err)
	}
	client := twitch.NewClient(cfg.ChannelName, cfg.OAuthToken)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if ratingStatus == 1 {
			checkRatingMsg(message)
		}

		if checkAccess(message) {
			checkCommand(message, client)
		}

	})

	client.Join(cfg.ChannelName)

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}
