package analysis

import "github.com/auho/go-toolkit-mysql/schema"

type Column struct {
	Column   schema.Column
	Amount   int
	Distinct int
	Empty    int
	Null     int
}
