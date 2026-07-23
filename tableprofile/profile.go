package tableprofile

import (
	simpledb "github.com/auho/go-simple-db/v3"
)

// Source identifies a MySQL table to analyse. It bundles the table name with
// the database connection used to query it.
type Source struct {
	Name string // table name
	DB   *simpledb.SimpleDB
}
