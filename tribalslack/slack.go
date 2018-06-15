package tribalslack

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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

func (slackConf SlackConfiguration) PostChallengeResponse(challenge string) error {
	slackClient := http.Client{}
	form := url.Values{}
	form.Add("token", slackConf.Token)
	form.Add("challenge", challenge)
	form.Add("type", "url_verification")
	req, err := http.NewRequest("POST", slackConf.SlackAPI, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = slackClient.Do(req)
	return err
}

func (sl SlackConfiguration)CheckMessageForChallengeAndRespond(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		if r.Body == nil {
			return errors.New("Please send a request body")
		}
		err := r.ParseForm()
		if err != nil {
			return err
		}
		challenge := r.Form.Get("challenge")
		token := r.Form.Get("token")
		if challenge != "" && token == sl.Token {
			sl.PostChallengeResponse(challenge)
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
