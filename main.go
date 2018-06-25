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
	"io/ioutil"
	"github.com/nlopes/slack/slackevents"
	"encoding/json"
)

type TribalServer struct {
	bot *bot.Bot
}

func main() {
	flag.Parse()
	api := slack.New("d7bffc23080740059982be80f194773e")
	users, err := api.GetUsers()
	fmt.Println(err)
	fmt.Println(users)
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
	http.HandleFunc("/challenge", s.ChallengePostHandler)
	http.HandleFunc("/request", s.RequestPostHandler)
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
	if r.Method == "POST" {
		if r.Body == nil {
			w.Write([]byte("Please Send a Body in your http Request"))
		}
		defer r.Body.Close()
		// TODO: determine criteria of message we care about
		command, err := slack.SlashCommandParse(r)
		if err != nil {
			fmt.Println("error in slash command parse")
			fmt.Println(err)
		}
		fmt.Println(command)
		switch caseVal := tribalslack.ParseCommand(command.Text); caseVal {
		case "rate":
			fmt.Println("rate found")
			err = s.bot.InitiateRateQuery(command, w, r)
		case "score":
			fmt.Println("score found")
			err = s.bot.InitiateScoreQuery(command, w, r)
		default:
			fmt.Println(fmt.Errorf("Unknown slash command"))
		}
		if err != nil {
			s.bot.InitiateError(command)
			fmt.Println(fmt.Errorf("Bad Slash Command"))
		}
	} else {
		w.Write([]byte("Please Send a POST http request"))
	}
	return
}

func (s TribalServer) ChallengePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.Body == nil {
			w.Write([]byte("Please Send a Body in your http Request"))
		}
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		err = tribalslack.CheckMessageForChallengeAndRespond(w, body)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s TribalServer) RequestPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.Body == nil {
			w.Write([]byte("Please Send a Body in your http Request"))
		}
		if err := r.ParseForm(); err != nil {
			fmt.Println(err)
			return
		}
		form := r.PostForm.Get("payload")
		var c slackevents.MessageAction
		err := json.Unmarshal([]byte(form), &c)
		fmt.Println(err)
		fmt.Println(c)
		fmt.Println(c.Actions[0])
	}
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

