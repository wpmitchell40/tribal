package tribalslack

import (
	"net/http"
	"strings"
	"time"
	"encoding/json"
	"fmt"
)
var Token = ""

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

type RateQueryFields struct {
	UserTakingQuery string
	UserBeingEvaluated string
}

type TribalQuery struct {
	Text string `json:"text"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Text string `json:"text"`
	Fallback string `json:"fallback"`
	CallbackId string `json:"callback_id"`
	Color string `json:"color"`
	AttachmentType string `json:"attachment_type"`
	Actions []SlackAction `json:"actions"`
}

type SlackAction struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Type string `json:"type"`
	Value string `json:"value"`
}

func PostChallengeResponse(w http.ResponseWriter, challenge string) error {
	_, err := w.Write([]byte(challenge))
	return err
}

func CheckMessageForChallengeAndRespond(w http.ResponseWriter, body []byte) error {

		var c Challenge
		err := json.Unmarshal(body, &c)
		if err != nil {
			fmt.Println("No Challenge Detected")
			return nil
		}
		if c.Challenge != "" && c.Token == Token {
			err = PostChallengeResponse(w, c.Challenge)
			if err != nil {
				fmt.Println(err)
				return err
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

func CreateTribalQuery() (TribalQuery){
	ActionYes := SlackAction{
		Name:"tribalresponse",
		Text:"Yes",
		Type:"button",
		Value:"Yes",
	}
	ActionNo := SlackAction{
		Name:"tribalresponse",
		Text:"No",
		Type:"button",
		Value:"No",
	}
	ActionNA := SlackAction{
		Name:"tribalresponse",
		Text:"N/A",
		Type:"button",
		Value:"N/A",
	}
	actions := []SlackAction{ActionYes,ActionNo,ActionNA}
	Attachment := SlackAttachment{
		Text:"If you do not feel you have interacted enough to answer, please choose N/A",
		Fallback:"An Error occurred, you were unable to rate the user",
		CallbackId:"tribal_response",
		Color:"#3AA3E3",
		AttachmentType:"default",
		Actions:actions,
	}
	Attachments := []SlackAttachment{Attachment}
	return TribalQuery{
		Text:"",
		Attachments:Attachments,
	}
}