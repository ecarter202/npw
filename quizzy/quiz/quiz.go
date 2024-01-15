package quiz

import (
	"fmt"
	"math/rand"
	"quizzy/models"
	"slices"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	correct = color.New(color.FgGreen)
	wrong   = color.New(color.FgRed)
)

type (
	Quiz struct {
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
func New(questions []*models.Question, opts ...OptsFunc) *Quiz {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	return &Quiz{
		Opts:      o,
		Questions: questions,
	}
}

func (quiz *Quiz) Start() {

	if quiz.IsRandomized {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(quiz.Questions), func(i, j int) {
			quiz.Questions[i], quiz.Questions[j] = quiz.Questions[j], quiz.Questions[i]
		})
	}

	for i, q := range quiz.Questions {
		if quiz.Length > 0 && i >= quiz.Length {
			break
		}

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
		quiz.total++ // add to possible points for this quiz

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
			correct.Println("Correct!")
			qScore := questionScore(userAnswers, correctAnswers)
			quiz.score += qScore // add to quiz score
		} else {
			answerLetters := mapKeys(correctAnswers)
			slices.Sort(answerLetters)
			wrong.Printf("Correct Answer: [%s]\n", strings.Join(answerLetters, ", "))
		}
	}

	fmt.Println("---------------------------------------------------------------")
	fmt.Printf("You got %d / %d correct which is %.2f%% \n", int(quiz.score), int(quiz.total), quiz.score/quiz.total*100)
	fmt.Println("---------------------------------------------------------------")
}

func questionScore(userAnswers, correctAnswers map[string]bool) (score float64) {
	l := float64(len(correctAnswers))
	for c := range correctAnswers {
		if userAnswers[c] != false {
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
