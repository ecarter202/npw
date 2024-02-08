package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"quizzy/config"
	"quizzy/models"
	"quizzy/shared"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	CSS_QUESTION_BOX  = "#adminForm > div.questions > div"
	CSS_QUESTION_TEXT = "div.panel-heading"
	CSS_ANSWER_BOX    = "ul > li"
)

func scrape(filesDir string) (cachedQuestions []*models.Question, err error) {
	cachedQuestions = []*models.Question{}
	seen := map[string]bool{}

	files, err := os.ReadDir(filesDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		filepath := fmt.Sprintf("%s/%s", config.HTML_FILES_DIR, f.Name())

		file, err := os.Open(filepath)
		if err != nil {
			return nil, fmt.Errorf("reading file [%s] [ERR: %s]\n", f, err)
		}
		defer file.Close()

		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			return nil, fmt.Errorf("document from reader (file) [%s] [ERR: %s]\n", f, err)
		}

		// iterate over the questions in HTML
		doc.Find(CSS_QUESTION_BOX).Each(func(i int, s *goquery.Selection) {
			qText := s.Find(CSS_QUESTION_TEXT).Text()
			qText = shared.SanitizeText(qText)

			if !seen[qText] {
				// create question

				h := md5.New()
				io.WriteString(h, qText)
				id := fmt.Sprintf("%x", h.Sum(nil))

				question := &models.Question{
					ID:      id,
					Text:    qText,
					Answers: []*models.Answer{},
				}

				// iterate over the answers for this question
				s.Find(CSS_ANSWER_BOX).Each(func(ii int, ss *goquery.Selection) {
					aText := ss.Text()
					answer := &models.Answer{
						Text: shared.SanitizeText(aText),
						// Letter: letterFromIndex(ii), // delete me
					}

					if strings.Contains(aText, "Missed)") {
						answer.IsCorrect = true
					}

					question.Answers = append(question.Answers, answer)
				})

				if len(question.Answers) > 0 {
					cachedQuestions = append(cachedQuestions, question)
				}
			}

		})
	}

	return cachedQuestions, nil
}

func letterFromIndex(i int) string {
	var alphabet = "abcdefghijklmnopqrstuvwxyz"
	return string(alphabet[i])
}
