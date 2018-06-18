package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"tribal_bot/bot"
	"tribal_bot/tribalslack"
	"tribal_bot/storage"

	"github.com/nlopes/slack"
)

type TribalServer struct {
	bot *bot.Bot
}

func main() {
	flag.Parse()
	api := slack.New("9uAyrqJby8XMCe8oM6UiWEfk")
	users, err := api.GetUsers()
	if err != nil {
		fmt.Println(err)
	}
	b, err := NewBot(*api, users)
	if err != nil {
		fmt.Println(err)
	}
	s := TribalServer{
		bot: b,

	}
	err = s.RunMetricsBot()
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func (s TribalServer) RunMetricsBot() (err error) {
	// TODO get an org list of users etc.
	// TODO establish a controller
	addr, err := determineListenAddress()
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/", s.SlashPostHandler)
	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
	return nil
}

func NewBot(client slack.Client, users []slack.User) (*bot.Bot, error) {
	tribalDB := storage.NewDB()
	helpbot := bot.NewBot(tribalDB, client, users)
	return helpbot, nil
}

func (s TribalServer) SlashPostHandler(w http.ResponseWriter, r *http.Request) {
	err := tribalslack.CheckMessageForChallengeAndRespond(w, r)
	if err != nil {
		fmt.Println(err)
	}
	// TODO: determine criteria of message we care about
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(command)
	switch caseVal := tribalslack.ParseCommand(command.Text); caseVal {
	case "rate":
		fmt.Println("rate found")
		err = s.bot.InitiateRateQuery(command)
	case "score":
		fmt.Println("score found")
		err = s.bot.InitiateScoreQuery(command)
	default:
		fmt.Println(fmt.Errorf("Unknown slash command"))
	}
	if err != nil {
		s.bot.InitiateError(command)
		fmt.Println(fmt.Errorf("Bad Slash Command"))
	}
	return
}

/*
func (s TribalServer) ControllerHandler(w http.ResponseWriter, r *http.Request) {
	slackMetricQuestion, err := s.slc.CreateSlackMessage(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.slc.PostNewMetric(slackMetricQuestion)
	s.bot.DataStorage.StoreMetricQuestionData(slackMetricQuestion)
	fmt.Println("Unable to log response")
}
*/

