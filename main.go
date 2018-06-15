package main

import (
	"encoding/json"
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

var (
	sl_conf      = flag.String("slconf", "/tribalslack/slackconf.json", "The location of the tribalslack configuration JSON file")
	storage_conf = flag.String("storageconf", "/storage/storageconf.json", "The location of the database configuration JSON file")
)

type TribalServer struct {
	bot *bot.Bot
	slconf *tribalslack.SlackConfiguration
}

func main() {
	flag.Parse()
	sl_conf, err := ConfigureSlack()
	if err != nil {
		log.Fatal(err)
	}
	api := slack.New(sl_conf.Token)
	users, err := api.GetUsers()
	if err != nil {
		log.Fatal(err)
	}
	b, err := NewBot(*api, users)
	if err != nil {
		log.Fatal(err)
	}
	s := TribalServer{
		bot: b,
		slconf:sl_conf,
	}
	err = s.RunMetricsBot()
}

func init() {
	flag.StringVar(sl_conf, "slc", "/tribalslack/slackconf.json", "The location of the tribalslack configuration JSON file")
	flag.StringVar(storage_conf, "dbc", "/storage/storageconf.json", "The location of the database configuration JSON file")
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
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/slash", s.SlashPostHandler)
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
	/*
	mux = http.NewServeMux()
	mux.HandleFunc("/controller", s.ControllerHandler)
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
	*/
	return nil
}

func NewBot(client slack.Client, users []slack.User) (*bot.Bot, error) {
	tribalDB := storage.NewDB()
	helpbot := bot.NewBot(tribalDB, client, users)
	return helpbot, nil
}

func (s TribalServer) SlashPostHandler(w http.ResponseWriter, r *http.Request) {
	err := s.slconf.CheckMessageForChallengeAndRespond(w, r)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: determine criteria of message we care about
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Fatal(err)
	}
	switch caseVal := tribalslack.ParseCommand(command.Text); caseVal {
	case "rate":
		err = s.bot.InitiateRateQuery(command)
	case "score":
		err = s.bot.InitiateScoreQuery(command)
	default:
		log.Fatal(fmt.Errorf("Unknown slash command"))
	}
	if err != nil {
		s.bot.InitiateError(command)
		log.Fatal(fmt.Errorf("Bad Slash Command"))
	}
	return
}

/*
func (s TribalServer) ControllerHandler(w http.ResponseWriter, r *http.Request) {
	slackMetricQuestion, err := s.slc.CreateSlackMessage(w, r)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = s.slc.PostNewMetric(slackMetricQuestion)
	s.bot.DataStorage.StoreMetricQuestionData(slackMetricQuestion)
	log.Println("Unable to log response")
}
*/

func ConfigureSlack() (*tribalslack.SlackConfiguration, error) {
	dir, _ := os.Getwd()
	path := dir + *sl_conf
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)
	var configuration tribalslack.SlackConfiguration
	if err := decoder.Decode(&configuration); err != nil {
		log.Fatal("error in decoder")
		return nil, err
	}
	return &configuration, nil
}

