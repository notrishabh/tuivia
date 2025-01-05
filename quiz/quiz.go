package quiz

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type QuizQuestion struct {
	Id                     int32               `json:"id"`
	Question               string              `json:"question"`
	Description            string              `json:"description"`
	Answers                map[string]string   `json:"answers"`
	MultipleCorrectAnswers string              `json:"multiple_correct_answers"`
	CorrectAnswers         map[string]string   `json:"correct_answers"`
	Explanation            string              `json:"explanation"`
	Tags                   []map[string]string `json:"tags"`
	Category               string              `json:"category"`
	Difficulty             string              `json:"difficulty"`
}

func Quiz() []QuizQuestion {
	apikey := os.Getenv("APIKEY")

	client := &http.Client{}

	baseurl := "https://quizapi.io/api/v1/questions"

	endpoint, err := url.Parse(baseurl)
	if err != nil {
		log.Fatal(err)
	}

	queryParams := url.Values{}
	queryParams.Set("apiKey", apikey)
	queryParams.Set("limit", "2")

	endpoint.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var questions []QuizQuestion
	err = json.Unmarshal(body, &questions)
	if err != nil {
		log.Fatal(err)
	}

	return questions
}
