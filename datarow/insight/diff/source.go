package diff

import (
	simpledb "github.com/auho/go-simple-db/v3"
)

type Source struct {
	Name string // table name
	DB   *simpledb.SimpleDB
}
