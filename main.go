package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/notrishabh/tuivia/quiz"
)

func main() {
	var (
		name string
	)
	var questions []quiz.QuizQuestion = quiz.Quiz()
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title(questions[0].Question).Options(
				huh.NewOption("Yushu boi", "yushu"),
				huh.NewOption("ippi boi", "ippi"),
				huh.NewOption("smol man", "smol"),
			).Value(&name),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}
