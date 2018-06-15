package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
	"tribal_bot/response"
	"tribal_bot/tribalslack"
	"errors"
)

type StorageDB struct {
	database *sql.DB
}

func NewDB() *StorageDB {
	db := StorageDB{}
	pdatabase, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	err = pdatabase.Ping()
	if err != nil {
		fmt.Println(err)
	}
	pdatabase.SetConnMaxLifetime(0)
	pdatabase.SetMaxIdleConns(0)
	db.database = pdatabase
	return &db
}

func (db StorageDB) StoreIntoResponses(response *response.TribalResponse) error {
	_, err := db.database.Exec(`
INSERT INTO responses (username, rating_user_username, response, response_time, organization_name, email_address)
VALUES ($1, $2, $3, $4, $5, $6)
`, response.Username, response.RespondingUserUsername, response.Response, time.Now(), response.OrgName, response.Email)
	return err
}

func (db StorageDB) GetUserScore(user string, queryFields *tribalslack.ScoreQueryFields) (float64, error) {
	rows, err := db.database.Query(`
SELECT response FROM responses 
WHERE responses.username = $1
AND responses.response_time >= $2; 
`, user, queryFields.Period)
	if err != nil {
		return 0, err
	}
	sum := 0.0
	count := 0
	for rows.Next() {
		var response string
		if err = rows.Scan(&response); err != nil {

			return 0,err
		}
		if response == "Y" {
			sum++
		} else if response == "N/A" {
			count--
		}
		count++
	}
	if count <= 10 {
		return 0, errors.New("Not enough users have rated this user")
	}
	tribalScore := sum/float64(count)
	return tribalScore, nil
}