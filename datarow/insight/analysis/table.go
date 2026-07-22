package analysis

import (
	"github.com/auho/go-toolkit-mysql/schema"
)

type Table struct {
	schema.Table
	Amount int
}
