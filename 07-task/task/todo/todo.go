package todo

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

var tasksBucket = []byte("tasks")

// Todo represents a task and its associated completion status
type Todo struct {
	Key         int    `json:"key"`
	Description string `json:"desc"`
	Completed   int64  `json:"completed"`
}

// FilterFn defines the signature for a predicate function to be used with GetFiltered()
type FilterFn func(t Todo) bool

// UpdateFn defines the signature for a function that updates a task
type UpdateFn func(t *Todo) error

// IsComplete indicates whether or not a Todo must still be accomplished
func (t Todo) IsComplete() bool {
	return t.Completed != 0
}

func (t Todo) String() string {

	var pre string

	if t.IsComplete() {
		pre = "[*]"
	} else {
		pre = "[ ]"
	}

	return fmt.Sprintf("%s %s", pre, t.Description)

}

// GetAll returns all Todos in the store
func GetAll() ([]Todo, error) {

	return GetFiltered(nil)

}

// GetFiltered returns Todos in the store that match a predicate
func GetFiltered(filter FilterFn) ([]Todo, error) {

	db, err := reader()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	results := []Todo{}

	err = db.View(func(tx *bolt.Tx) error {

		todos := []Todo{}

		c := tx.Bucket(tasksBucket).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var todo Todo

			err := json.Unmarshal(v, &todo)

			if err != nil {
				log.Println(err)
				continue
			}

			if filter == nil || filter(todo) {
				todos = append(todos, todo)
			}

		}

		results = make([]Todo, len(todos))

		copy(results, todos)

		return nil

	})

	if err != nil {
		return nil, err
	}

	return results, nil

}

// Create creates a new Todo and stores it
func Create(desc string) ([]Todo, error) {

	db, err := writer()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	result := make([]Todo, 1)

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(tasksBucket)

		keyInt64, err := b.NextSequence()

		if err != nil {
			return err
		}

		key := int(keyInt64)

		todo := Todo{
			Key:         key,
			Description: desc,
			Completed:   0,
		}

		bytes, err := json.Marshal(todo)

		if err != nil {
			return err
		}

		err = b.Put(itob(key), bytes)

		if err != nil {
			return err
		}

		result[0] = todo

		return nil

	})

	if err != nil {
		return nil, err
	}

	return result, nil

}

// Get returns the Todo specified by a given key
func Get(key int) ([]Todo, error) {

	return Update(key, nil)

}

// Update performs an update action on a Task specified by a given key
func Update(key int, updateFn UpdateFn) ([]Todo, error) {

	db, err := writer()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	result := make([]Todo, 1)

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(tasksBucket)

		c := b.Cursor()

		k, v := c.Seek(itob(key))

		if k == nil {
			return fmt.Errorf("no task with given key '%d'", key)
		}

		var todo Todo

		err := json.Unmarshal(v, &todo)

		if err != nil {
			return fmt.Errorf("failed to unmarshall task with given key '%d': %s", key, err)
		}

		if updateFn != nil {
			if err = updateFn(&todo); err != nil {
				return fmt.Errorf("failed to update task with given key '%d': %s", key, err)
			}
			if v, err = json.Marshal(todo); err != nil {
				return fmt.Errorf("failed to marshal updated task with given key '%d': %s", key, err)
			}
			if err = b.Put(itob(key), v); err != nil {
				return fmt.Errorf("failed to write back updated task with given key '%d': %s", key, err)
			}
		}

		result[0] = todo

		return nil

	})

	if err != nil {
		return nil, err
	}

	return result, nil

}

// Delete removes a task from the task store, given its key
func Delete(key int) error {

	db, err := writer()

	if err != nil {
		return err
	}

	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(tasksBucket)

		c := b.Cursor()

		k, _ := c.Seek(itob(key))

		if k == nil {
			return fmt.Errorf("no task with key %d", key)
		}

		return b.Delete(k)

	})

}

// Init creates the task store (unless one exists) and guarantees buckets
func Init() error {

	// Initialize tha Bolt database to store our tasks

	db, err := writer()

	if err != nil {
		return err
	}

	defer db.Close()

	// Guarantee that the tasks bucket exists

	err = db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(tasksBucket)

		if err != nil {
			return fmt.Errorf("called CreateBucketIfNotExists() with name '%s'", tasksBucket)
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil

}

func reader() (*bolt.DB, error) {

	db, err := bolt.Open("tasks.db", 0666, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})

	if err != nil {
		return nil, err
	}

	return db, nil

}

func writer() (*bolt.DB, error) {

	db, err := bolt.Open("tasks.db", 0666, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return nil, err
	}

	return db, nil

}

func itob(value int) []byte {

	valueBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(valueBytes, uint64(value))

	return valueBytes

}
