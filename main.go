package main

import (
	"flag"
	"log"
	"quizzy/config"
	"quizzy/database"
	"quizzy/models"
	"quizzy/quiz"
)

var (
	_missed        bool
	_questionCount = 25

	err error
)

func main() {
	flag.BoolVar(&_missed, "m", _missed, "Get only the previously missed questions.")
	flag.IntVar(&_questionCount, "n", _questionCount, "Number of questions in your quiz.")
	flag.Parse()

	var (
		questions []*models.Question
		db        *database.Database
	)

	db, err = database.New(config.DB_NAME)
	if err != nil {
		log.Fatalf("creating database [ERR: %s]", err)
	}

	if _missed {
		questions, err = db.AllMissed()
		if err != nil {
			log.Fatal(err)
		}
		_questionCount = len(questions)
	} else {
		questions, err = scrape(config.HTML_FILES_DIR)
		if err != nil {
			log.Fatal(err)
		}
	}
	db.Close()

	if len(questions) == 0 {
		log.Fatal("no questions found...")
	}

	quiz.New(questions,
		quiz.OptRandomize,
		quiz.OptSetLength(_questionCount),
	).Start()
}
