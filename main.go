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
		name string
	)
	var questions []quiz.QuizQuestion = quiz.Quiz()
	for i, v := range questions[0].Answers {
		fmt.Println(i, v)
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title(questions[0].Question).Options(
				huh.NewOption("Yushu boi", "yushu"),
				huh.NewOption("ippi boi", "ippi"),
				huh.NewOption("smol man", "smol"),
			).Value(&name),
		),
	)

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}
