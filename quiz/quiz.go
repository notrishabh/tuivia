package quiz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
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
	AnswersArray           []string
	CorrectAnswer          string
}

type quesOrCate interface {
	~[]QuizQuestion | ~[]Category
}

func fetchApi[T quesOrCate](ctx context.Context, url string, results chan<- string, resData chan<- T) {
	go func() {
		defer close(results)
		defer close(resData)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			results <- fmt.Sprintf("Error creating req for %s: %s", url, err.Error())
			return
		}

		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			results <- fmt.Sprintf("Error making req to %s: %s", url, err.Error())
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- fmt.Sprintf("Error reading response body from %s: %s", url, err.Error())
			return
		}
		var data T
		err = json.Unmarshal(body, &data)
		if err != nil {
			results <- fmt.Sprintf("Error unmarshalling response body from %s: %s", url, err.Error())
			return
		}

		resData <- data

	}()
}

func Quiz(category string) ([]QuizQuestion, error) {
	var wg sync.WaitGroup

	apikey := os.Getenv("APIKEY")

	baseurl := "https://quizapi.io/api/v1/questions"

	endpoint, err := url.Parse(baseurl)
	if err != nil {
		log.Fatal(err)
	}

	queryParams := url.Values{}
	queryParams.Set("apiKey", apikey)
	queryParams.Set("limit", "5")
	if category != "all" {
		queryParams.Set("category", category)
	}

	endpoint.RawQuery = queryParams.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res := make(chan string, 1)
	questionsChan := make(chan []QuizQuestion, 1)

	fetchApi(ctx, endpoint.String(), res, questionsChan)

	select {
	case err := <-res:
		return nil, fmt.Errorf(err)
	case <-ctx.Done():
		return nil, fmt.Errorf("Task timed out")
	case ques := <-questionsChan:
		if ques == nil {
			return nil, fmt.Errorf("No qestions received from api")
		}

		out := make(chan QuizQuestion)
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer close(out)
			for i := range ques {
				for _, v := range ques[i].Answers {
					if v != "" {
						ques[i].AnswersArray = append(ques[i].AnswersArray, v)
					}
				}

				for k, v := range ques[i].CorrectAnswers {
					if v == "true" {
						ques[i].CorrectAnswer = ques[i].AnswersArray[mapCorrectAns(k)]
					}

				}
			}
		}()
		wg.Wait()
		return ques, nil
	}
}

func mapCorrectAns(key string) int {
	switch key {
	case "answer_a_correct":
		return 0
	case "answer_b_correct":
		return 1
	case "answer_c_correct":
		return 2
	case "answer_d_correct":
		return 3
	case "answer_e_correct":
		return 4
	case "answer_f_correct":
		return 5
	default:
		return 0
	}
}
