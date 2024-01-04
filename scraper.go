package main

import (
	"fmt"
	"nwp/config"
	"nwp/models"
	"nwp/shared"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	CSS_QUESTION_BOX  = "#adminForm > div.questions > div"
	CSS_QUESTION_TEXT = "div.panel-heading"
	CSS_ANSWER_BOX    = "ul > li"
)

func scrape(filesDir string) error {
	cachedQuestions = []*models.Question{}
	seen := map[string]bool{}

	files, err := os.ReadDir(filesDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		filepath := fmt.Sprintf("%s/%s", config.HTML_FILES_DIR, f.Name())

		file, err := os.Open(filepath)
		if err != nil {
			return fmt.Errorf("reading file [%s] [ERR: %s]\n", f, err)
		}
		defer file.Close()

		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			return fmt.Errorf("document from reader (file) [%s] [ERR: %s]\n", f, err)
		}

		// iterate over the questions in HTML
		doc.Find(CSS_QUESTION_BOX).Each(func(i int, s *goquery.Selection) {
			qText := s.Find(CSS_QUESTION_TEXT).Text()
			qText = shared.SanitizeText(qText)

			if !seen[qText] {
				// create question

				question := &models.Question{
					Text:    qText,
					Answers: []*models.Answer{},
				}

				// iterate over the answers for this question
				s.Find(CSS_ANSWER_BOX).Each(func(ii int, ss *goquery.Selection) {
					answer := &models.Answer{
						Text:   shared.SanitizeText(ss.Text()),
						Letter: letterFromIndex(ii),
					}

					if strings.Contains(ss.Text(), "( Missed)") {
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

	// print out all questions as JSON
	// b, err := json.MarshalIndent(cachedQuestions, "", "    ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(b))

	return nil
}

func letterFromIndex(i int) string {
	var alphabet = "abcdefghijklmnopqrstuvwxyz"
	return string(alphabet[i])
}
