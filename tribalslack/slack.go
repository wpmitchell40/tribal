package tribalslack

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
	"fmt"
)
var SlackAPI ="https://api.groupme.com/v3/bots/post"
var Token = "9uAyrqJby8XMCe8oM6UiWEfk"

type SlackConfiguration struct {
	SlackAPI     string
	Token        string
	ClientID     string
	ClientSecret string
	AppID        string
}

type ScoreQueryFields struct {
	User     string
	Period time.Time
	Report bool
}

func PostChallengeResponse(challenge string) error {
	slackClient := http.Client{}
	form := url.Values{}
	form.Add("token", Token)
	form.Add("challenge", challenge)
	form.Add("type", "url_verification")
	req, err := http.NewRequest("POST", SlackAPI, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = slackClient.Do(req)
	return err
}

func CheckMessageForChallengeAndRespond(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		if r.Body == nil {
			return errors.New("Please send a request body")
		}
		err := r.ParseForm()
		if err != nil {
			return err
		}
		fmt.Println(r.Form)
		challenge := r.Form.Get("challenge")
		fmt.Println(challenge)
		token := r.Form.Get("token")
		fmt.Println(token)
		if challenge != "" && token == Token {
			PostChallengeResponse(challenge)
			return nil
		}
	}
	return nil
}

func ParseCommand(command string)(string){
	if strings.Contains(command, "rate") {
		return "rate"
	} else if strings.Contains(command, "score") {
		return "score"
	}
	return ""
}
