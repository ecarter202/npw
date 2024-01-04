package main

import (
	"fmt"
	"nwp/models"
	"strings"
)

var (
	cachedQuestions []*models.Question
)

func startQuiz(questions []*models.Question) {
	for _, q := range questions {
		var (
			userInput      string
			userAnswers    = map[string]bool{}
			correctAnswers = map[string]bool{}
		)

		for _, a := range q.Answers {
			if a.IsCorrect {
				correctAnswers[a.Letter] = true
			}
		}

		fmt.Println(q.Text)
		for _, a := range q.Answers {
			fmt.Printf("%s) %s\n", a.Letter, a.Text)
		}
		fmt.Scanln(&userInput)

		x := strings.Split(userInput, ",")
		for _, str := range x {
			userAnswerLetter := strings.TrimSpace(str)
			userAnswers[userAnswerLetter] = true
		}

		if answeredCorrectly(userAnswers, correctAnswers) {
			fmt.Println("GOOD JOB!")
		} else {
			fmt.Println("YOU SUCK!!")
		}
	}
}

func answeredCorrectly(users, correct map[string]bool) bool {
	fmt.Println("Checking your answers of ", users)
	fmt.Println("Against correct answers of ", correct)
	for c := range correct {
		if users[c] == false {
			return false
		}
	}

	return true
}
