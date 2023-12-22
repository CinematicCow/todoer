package db

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

const (
	dbName = "todoer.db"
	dbFile = ".todoer"
)

type Todo struct {
	ID        uint64
	Text      string
	Completed bool
}

// initializes the database and creates the necessary bucket for todos.
func init() {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not get home dir: ", err)
	}
	dbDir := filepath.Join(home, dbFile)

	err = os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		log.Fatal("Could not create db dir: ", err)
	}
	dbPath := filepath.Join(dbDir, dbName)

	df, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal("Could not open db: ", err)
	}
	db = df

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("todos"))
		return err
	})
	if err != nil {
		log.Fatal("Could not create bucket: ", err)
	}

}

// adds a todo to the database.
//
// It takes a string parameter `todo` representing the todo to be added.
// It returns an error indicating any issues encountered during the process.
func AddToDB(todo string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("todos"))
		id, _ := b.NextSequence()

		// Create a new todo instance
		t := Todo{
			ID:        id,
			Text:      todo,
			Completed: false,
		}

		// Marshal the todo
		d, err := json.Marshal(t)
		if err != nil {
			return err
		}

		// Store the marshalled todo in the bucket
		return b.Put([]byte(strconv.FormatUint(id, 10)), d)
	})
}

// retrieves all the todos from the database.
//
// It returns a slice of Todo structs and an error.
// It returns an error indicating any issues encountered during the process.
func GetAllTodos() ([]Todo, error) {
	var todos []Todo
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("todos"))
		return b.ForEach(func(k, v []byte) error {
			var t Todo
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			todos = append(todos, t)
			return nil
		})
	})
	return todos, err
}

// marks a todo item as completed in the database.
//
// It takes an unsigned 64-bit integer `id` as a parameter representing the ID of the todo item to be marked as completed.
// It returns an error indicating any issues encountered during the process.
func MarkAsCompleted(id uint64) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("todos"))
		v := b.Get([]byte(strconv.FormatUint(id, 10)))
		var t Todo
		if err := json.Unmarshal(v, &t); err != nil {
			return err
		}
		t.Completed = true
		d, err := json.Marshal(t)
		if err != nil {
			return err
		}

		return b.Put([]byte(strconv.FormatUint(id, 10)), d)
	})
}

// deletes all completed todos from the database.
//
// It does this by iterating over each key-value pair in the "todos" bucket of the database.
// For each pair, it unmarshals the value into a Todo struct and checks if the "Completed" field is true.
// If it is, the pair is deleted from the bucket.
//
// Returns an error if there was an issue with the database transaction or unmarshaling the value.
func DeleteCompletedTodos() (int, error) {
	var count int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("todos"))
		return b.ForEach(func(k, v []byte) error {
			var t Todo
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			if t.Completed {
				err := b.Delete(k)
				if err != nil {
					return err
				}
				count++
			}
			return nil
		})
	})
	return count, err
}
