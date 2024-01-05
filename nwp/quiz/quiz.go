package quiz

import (
	"fmt"
	"math/rand"
	"nwp/models"
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
				quiz.total++ // add to possible points for this quiz
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
			correct.Println("Correct!")
			quiz.score++ // add to quiz score
		} else {
			wrong.Printf("Correct Answer: [%s]\n", strings.Join(mapKeys(correctAnswers), ", "))
		}
	}

	fmt.Println("---------------------------------------------------------------")
	fmt.Printf("You got %d / %d correct which is %.2f%% \n", int(quiz.score), int(quiz.total), quiz.score/quiz.total*100)
	fmt.Println("---------------------------------------------------------------")
}

func answeredCorrectly(users, correct map[string]bool) bool {
	for c := range correct {
		if users[c] == false {
			return false
		}
	}

	return true
}

func mapKeys(m map[string]bool) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}

	return
}
