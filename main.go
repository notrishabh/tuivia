package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/joho/godotenv"
	"github.com/notrishabh/tuivia/quiz"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		name     string
		question string
	)
	var questions []quiz.QuizQuestion = quiz.Quiz()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title(questions[0].Question).Options(
				huh.NewOption("Yushu boi", "yushu"),
				huh.NewOption("ippi boi", "ippi"),
				huh.NewOption("smol man", "smol"),
			).Value(&name),
			huh.NewSelect[string]().Value(&question).Title(questions[0].Question).OptionsFunc(func() []huh.Option[string] {
				var opts []string
				for _, v := range questions[0].Answers {
					if v != "" {
						opts = append(opts, v)
					}
				}
				return huh.NewOptions(opts...)
			}, &question),
		),
	)

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}
