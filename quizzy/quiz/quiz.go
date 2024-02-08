package quiz

import (
	"fmt"
	"log"
	"math/rand"
	"quizzy/config"
	"quizzy/database"
	"quizzy/models"
	"slices"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	blue  = color.New(color.FgHiBlue)
	green = color.New(color.FgGreen)
	red   = color.New(color.FgRed)
)

type (
	Quiz struct {
		*database.Database
		*models.Subject

		Opts
		Questions []*models.Question

		total float64
		score float64
	}

	Opts struct {
		IsRandomized bool
		Length       int // num of questions to ask
	}

	OptsFunc func(*Opts)
)

func defaultOpts() Opts {
	// default is empty for now
	// using type's zero values
	return Opts{}
}

func OptRandomize(opts *Opts) {
	opts.IsRandomized = true
}

func OptSetLength(i int) OptsFunc {
	return func(opts *Opts) {
		opts.Length = i
	}
}

// New creates a new quiz
func New(questionPool []*models.Question, opts ...OptsFunc) *Quiz {
	db, err := database.New(config.DB_NAME)
	if err != nil {
		log.Fatalf("creating database [ERR: %s]", err)
	}

	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	quiz := &Quiz{
		Database: db,
		Opts:     o,
	}

	if quiz.IsRandomized {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(questionPool), func(i, j int) {
			questionPool[i], questionPool[j] = questionPool[j], questionPool[i]
		})
		for _, q := range questionPool {
			r.Shuffle(len(q.Answers), func(i, j int) {
				q.Answers[i], q.Answers[j] = q.Answers[j], q.Answers[i]
			})
		}
	}

	questions := questionPool[:quiz.Length]

	quiz.Questions = questions
	quiz.total = float64(len(questions))

	return quiz
}

func (quiz *Quiz) Start() {
	blue.Printf("%d questions. Go!\n", len(quiz.Questions))
	fmt.Println()

	for i, q := range quiz.Questions {
		if quiz.Length > 0 && i >= quiz.Length {
			break
		}

		var (
			userInput      string
			userAnswers    = map[string]bool{}
			correctAnswers = map[string]bool{}
		)

		for ii, a := range q.Answers {
			a.Letter = letterFromIndex(ii)
			if a.IsCorrect {
				correctAnswers[a.Letter] = true
			}
		}

		fmt.Printf("%d) %s\n", i+1, q.Text)
		for _, a := range q.Answers {
			fmt.Printf("	%s) %s\n", a.Letter, a.Text)
		}
		fmt.Scanln(&userInput)

		x := strings.Split(userInput, ",")
		for _, str := range x {
			userAnswerLetter := strings.TrimSpace(str)
			userAnswers[userAnswerLetter] = true
		}

		if answeredCorrectly(userAnswers, correctAnswers) {
			green.Println("Correct!")
			quiz.DelMissed(q.ID)
		} else {
			answerLetters := mapKeys(correctAnswers)
			slices.Sort(answerLetters)
			red.Printf("Correct Answer: [%s]\n", strings.Join(answerLetters, ", "))

			err := quiz.AddMissed(q)
			if err != nil {
				log.Fatalf("adding missed question \n%s\n [ERR: %s]", q.JSON(), err)
			}
		}

		qScore := questionScore(userAnswers, correctAnswers)
		quiz.score += qScore // add to quiz score
	}

	fmt.Println("---------------------------------------------------------------")
	fmt.Printf("You got %.2f / %d correct which is %.2f%% \n", quiz.score, int(quiz.total), quiz.score/quiz.total*100)
	fmt.Println("---------------------------------------------------------------")
}

func questionScore(userAnswers, correctAnswers map[string]bool) (score float64) {
	l := float64(len(correctAnswers))
	for answer := range correctAnswers {
		if userAnswers[answer] != false {
			score += 100 / (100 * l)
		}
	}

	return score
}

func answeredCorrectly(userAnswers, correctAnswers map[string]bool) bool {
	// 1 being max points per question
	return questionScore(userAnswers, correctAnswers) == 1
}

func mapKeys(m map[string]bool) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}

	return
}

func letterFromIndex(i int) string {
	var alphabet = "abcdefghijklmnopqrstuvwxyz"
	return string(alphabet[i])
}
