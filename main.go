package main

import (
	"flag"
	"fmt"
	"log"
	"nwp/config"
	"nwp/models"
	"nwp/quiz"

	bolt "go.etcd.io/bbolt"
)

var (
	setup bool

	db  *bolt.DB
	err error
)

func main() {
	flag.BoolVar(&setup, "init", setup, "Initialize setup? This populates the DB with questions.")
	flag.Parse()

	var questions []*models.Question
	if setup {
		// maybe use a db? not sure yet...
		// err = initDB()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer db.Close()

		questions, err = scrape(config.HTML_FILES_DIR)
		if err != nil {
			log.Fatal(err)
		}
	}

	quiz.New(questions,
		quiz.OptRandomize,
		quiz.OptSetLength(5),
	).Start()
}

func initDB() error {
	db, err = bolt.Open(config.DB_NAME, 0600, nil)
	if err != nil {
		log.Fatalf("opening DB [%s] [ERR: %s]\n", config.DB_NAME, err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(config.DB_Q_BUCKET))
		if err != nil {
			return fmt.Errorf("creating bucket [%s]: %s", config.DB_Q_BUCKET, err)
		}
		return nil
	})

	return nil
}
