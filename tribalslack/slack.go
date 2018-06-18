package tribalslack

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"
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

type Challenge struct {
	Type string
	Token string
	Challenge string
}

type ChallengeResonse struct {
	Challenge string
}

type ScoreQueryFields struct {
	User     string
	Period time.Time
	Report bool
}

func PostChallengeResponse(w http.ResponseWriter, challenge string) error {
	_, err := w.Write([]byte(challenge))
	return err
}

func CheckMessageForChallengeAndRespond(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		if r.Body == nil {
			return errors.New("Please send a request body")
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		var c Challenge
		err = json.Unmarshal(body, &c)
		if err != nil {
			fmt.Println("No Challenge Detected")
			return nil
		}
		defer r.Body.Close()
		if c.Challenge != "" && c.Token == Token {
			err = PostChallengeResponse(w, c.Challenge)
			if err != nil {
				fmt.Println(err)
				return err
			}
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
