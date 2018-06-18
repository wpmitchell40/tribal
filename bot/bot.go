package bot

import (
	"tribal_bot/response"
	"tribal_bot/storage"
	"tribal_bot/tribalslack"
	"strings"
	"github.com/nlopes/slack"
	"errors"
	"time"
	"strconv"
	"fmt"
)

type Bot struct {
	DataStorage *storage.StorageDB
	SlackAPI slack.Client
	Users []slack.User
}

func (b *Bot) LogInteraction(response *response.TribalResponse) (err error) {
	return b.DataStorage.StoreIntoResponses(response)

}

func NewBot(storage *storage.StorageDB, client slack.Client, users []slack.User) *Bot {
	return &Bot{
		DataStorage: storage,
		SlackAPI:client,
		Users:users,
	}
}

func (b *Bot) InitiateError(command slack.SlashCommand) {
	params := slack.PostMessageParameters{}
	b.SlackAPI.PostMessage(command.UserID,
		"An error occurred with your most recent Tribal command, please check your syntax and try again", params)
}

func (b *Bot) InitiateRateQuery(command slack.SlashCommand) error {
	return nil
}

func (b *Bot) InitiateScoreQuery(command slack.SlashCommand) error {
	 queryFields, err := b.parseCommandForQueryFields(command.Text)
	 if err != nil {
	 	return err
	 }
	 fmt.Println("QueryFields:")
	 fmt.Println(queryFields)
	 validUser := false
	 index := -1
	 for i, v := range b.Users {
	 	if strings.Contains(queryFields.User,v.ID) {
	 		validUser = true
	 		index = i
		}
	 }
	 if validUser {
	 	userToScore := b.Users[index].ID
	 	params := slack.PostMessageParameters{}
	 	score, err := b.DataStorage.GetUserScore(userToScore, queryFields)
	 	if err != nil {
			b.SlackAPI.PostMessage(command.UserID,"", params)
		}
		text := fmt.Sprintf("The TribalScore for that user in that timeframe is %f", score)
	 	b.SlackAPI.PostMessage(command.UserID, text, params)
	 }

	return nil
}

func (b *Bot) parseCommandForQueryFields(text string) (*tribalslack.ScoreQueryFields, error) {
	fields := strings.Fields(text)
	if len(fields) < 2 || len(fields) > 4 {
		return nil, errors.New("score slash command is invalid")
	}
	scoreFields := tribalslack.ScoreQueryFields{}
	scoreFields.User = fields[1]
	year, err := parseDuration("1y")
	if err != nil {
		return nil, err
	}
	if len (fields) == 2 {
		scoreFields.Report = false
		scoreFields.Period = *year
		return &scoreFields, nil
	}
	if len(fields) == 3 {
		if strings.Contains(text, "report"){
			scoreFields.Report = true
			scoreFields.Period = *year
			return &scoreFields, nil
		} else {
			scoreFields.Report = false
			maxTimeAway, err := parseDuration(fields[2])
			if err != nil {
				return nil, nil
			}
			scoreFields.Period = *maxTimeAway
			return &scoreFields, nil
		}
	} else{
		scoreFields.Report = true
		maxTimeAway, err := parseDuration(fields[2])
		if err != nil {
			return nil, nil
		}
		scoreFields.Period = *maxTimeAway
		return &scoreFields, nil
	}
	return nil, nil
}

func parseDuration(s string) (*time.Time, error) {
	maxTimeAway := time.Now()
	if idx := strings.IndexByte(s, 'y'); idx >= 0 {
		yearValue := s[:idx]
		s = s[idx+1:]
		y, err := strconv.Atoi(yearValue)
		if err != nil {return nil, err}
		hoursInYears := 24*365*y
		h := strconv.Itoa(hoursInYears)
		h = h + "h"
		dur, err := time.ParseDuration(h)
		if err != nil {return nil, err}
		maxTimeAway.Add(-dur)
	}
	if idx := strings.IndexByte(s, 'm'); idx >= 0 {
		monthValue := s[:idx]
		s = s[idx+1:]
		m, err := strconv.Atoi(monthValue)
		if err != nil {return nil, err}
		hoursInMonths := 24*30*m
		h := strconv.Itoa(hoursInMonths)
		h = h + "h"
		dur, err := time.ParseDuration(h)
		if err != nil {return nil, err}
		maxTimeAway.Add(-dur)
	}
	if idx := strings.IndexByte(s, 'w'); idx >= 0 {
		weekValue := s[:idx]
		s = s[idx+1:]
		w, err := strconv.Atoi(weekValue)
		if err != nil {return nil, err}
		hoursInWeeks := 24*7*w
		h := strconv.Itoa(hoursInWeeks)
		h = h + "h"
		dur, err := time.ParseDuration(h)
		if err != nil {return nil, err}
		maxTimeAway.Add(-dur)
	}
	if idx := strings.IndexByte(s, 'd'); idx >= 0 {
		dayValue := s[:idx]
		s = s[idx+1:]
		d, err := strconv.Atoi(dayValue)
		if err != nil {return nil, err}
		hoursInDays := 24*d
		h := strconv.Itoa(hoursInDays)
		h = h + "h"
		dur, err := time.ParseDuration(h)
		if err != nil {return nil, err}
		maxTimeAway.Add(-dur)
	}
	return &maxTimeAway, nil
}