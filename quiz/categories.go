package quiz

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
)

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetCategories() ([]Category, error) {
	apiKey := os.Getenv("APIKEY")

	baseUrl := "https://quizapi.io/api/v1/categories"

	endpoint, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatal(err)
	}

	queryParams := url.Values{}
	queryParams.Set("apiKey", apiKey)

	endpoint.RawQuery = queryParams.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	categories := make(chan []Category)
	res := make(chan string, 1)
	fetchApi(ctx, endpoint.String(), res, categories)

	select {
	case err := <-res:
		return nil, errors.New(err)
	case <-ctx.Done():
		return nil, fmt.Errorf("Task timed out")
	case cate := <-categories:
		if cate == nil {
			return nil, fmt.Errorf("No categories found.")
		}
		return cate, nil
	}

}
