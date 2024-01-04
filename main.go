package main

import (
	"flag"
	"fmt"
	"log"
	"nwp/config"

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

	if setup {
		// maybe use a db? not sure yet...
		// err = initDB()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer db.Close()

		err := scrape(config.HTML_FILES_DIR)
		if err != nil {
			log.Fatal(err)
		}
	}

	startQuiz(cachedQuestions)

	fmt.Println("DONE")
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
