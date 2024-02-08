package database

import (
	"encoding/json"
	"fmt"
	"log"
	"quizzy/config"
	"quizzy/models"

	bolt "go.etcd.io/bbolt"
)

type (
	Database struct {
		*bolt.DB
	}
)

func New(name string) (*Database, error) {
	db, err := bolt.Open(name, 0600, nil)
	if err != nil {
		log.Fatalf("opening DB [%s] [ERR: %s]\n", config.DB_NAME, err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(config.DB_MISSED_BUCKET))
		if err != nil {
			return fmt.Errorf("creating bucket [%s]: %s", config.DB_MISSED_BUCKET, err)
		}
		return nil
	})

	return &Database{db}, nil
}

func (db *Database) AllMissed() (questions []*models.Question, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte(config.DB_MISSED_BUCKET)).
			ForEach(func(k, v []byte) error {
				q := new(models.Question)
				err = json.Unmarshal(v, &q)
				if err != nil {
					return err
				}

				questions = append(questions, q)

				return nil
			})

		return err
	})

	return
}

func (db *Database) AddMissed(question *models.Question) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.DB_MISSED_BUCKET))
		err := bucket.Put([]byte(question.ID), []byte(question.JSON()))
		return err
	})

	return err
}

func (db *Database) DelMissed(questionID string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(config.DB_MISSED_BUCKET))
		err := bucket.Delete([]byte(questionID))
		return err
	})

	return err
}
