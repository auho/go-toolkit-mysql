package explore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/auho/go-toolkit-mysql/internal/testutil/mysql"
)

// ddl for the explore test table
const exploreTestDDL = "CREATE TABLE `explore_test` (" +
	"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
	"`i` int(11) NOT NULL DEFAULT '0'," +
	"`s` varchar(20) NOT NULL DEFAULT ''," +
	"`s_null` varchar(20) DEFAULT NULL," +
	"PRIMARY KEY (`id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

func TestExplorer_Analyze(t *testing.T) {
	simpleDB, gormDB, err := mysql.NewDB()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = simpleDB.Drop(ctx, "explore_test")
	})

	err = mysql.RecreateTable(gormDB, "explore_test", exploreTestDDL)
	if err != nil {
		t.Fatalf("recreate table: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// insert test data:
	//   i:        0, 1, 2, 3, 4, 4, 4   -> distinct=5, empty(=0)=1, null=0
	//   s:        "", 1, 2, 3, 4, 4, 4  -> distinct=5, empty(='')=1, null=0
	//   s_null:   nil, "", 1, 2, 3, 3, 3 -> distinct=4 (NULL excluded), empty(='')=1, null=1
	err = simpleDB.BulkInsertFromSliceSlice(ctx, "explore_test",
		[]string{"i", "s", "s_null"},
		[][]any{
			{0, "", nil},
			{1, "1", ""},
			{2, "2", "1"},
			{3, "3", "2"},
			{4, "4", "3"},
			{4, "4", "3"},
			{4, "4", "3"},
		},
		100,
	)
	if err != nil {
		t.Fatalf("insert data: %v", err)
	}

	e := New(simpleDB)
	result, err := e.Analyze(ctx, "explore_test")
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	// table-level checks
	if result.Table.Name != "explore_test" {
		t.Errorf("Table.Name = %q, want %q", result.Table.Name, "explore_test")
	}
	if result.Table.RowCount != 7 {
		t.Errorf("Table.RowCount = %d, want 7", result.Table.RowCount)
	}

	// field names should match column order
	wantFields := []string{"id", "i", "s", "s_null"}
	if len(result.FieldNames) != len(wantFields) {
		t.Fatalf("FieldNames len = %d, want %d", len(result.FieldNames), len(wantFields))
	}
	for k, fn := range wantFields {
		if result.FieldNames[k] != fn {
			t.Errorf("FieldNames[%d] = %q, want %q", k, result.FieldNames[k], fn)
		}
	}

	// per-column checks
	type colExpect struct {
		name     string
		rowCount int
		distinct int
		empty    int
		null     int
	}
	expects := []colExpect{
		{"id", 7, 7, 0, 0},
		{"i", 7, 5, 1, 0},
		{"s", 7, 5, 1, 0},
		{"s_null", 7, 4, 1, 1},
	}

	for _, exp := range expects {
		c, ok := result.Columns[exp.name]
		if !ok {
			t.Errorf("column %q not found in result", exp.name)
			continue
		}
		if c.RowCount != exp.rowCount {
			t.Errorf("column %q: RowCount = %d, want %d", exp.name, c.RowCount, exp.rowCount)
		}
		if c.Distinct != exp.distinct {
			t.Errorf("column %q: Distinct = %d, want %d", exp.name, c.Distinct, exp.distinct)
		}
		if c.Empty != exp.empty {
			t.Errorf("column %q: Empty = %d, want %d", exp.name, c.Empty, exp.empty)
		}
		if c.Null != exp.null {
			t.Errorf("column %q: Null = %d, want %d", exp.name, c.Null, exp.null)
		}
	}

	// ToStrings should produce output without error
	lines := result.ToStrings()
	if len(lines) == 0 {
		t.Error("ToStrings returned empty slice")
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}
