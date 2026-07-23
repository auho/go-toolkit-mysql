package tableprofile

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	simpledb "github.com/auho/go-simple-db/v3"
	"github.com/auho/go-toolkit-mysql/internal/testutil/mysql"
)

// shared database connection and test tables used by all tests in this package
var simpleDB *simpledb.SimpleDB

const diffDDL = "CREATE TABLE `diff` (" +
	"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
	"`i` int(11) NOT NULL DEFAULT '0'," +
	"`s` varchar(20) NOT NULL DEFAULT ''," +
	"`s_null` varchar(20) DEFAULT NULL," +
	"`d1` varchar(20) DEFAULT NULL," +
	"`中文字段1` varchar(20) DEFAULT NULL," +
	"PRIMARY KEY (`id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

const diffCopyDDL = "CREATE TABLE `diff_copy` (" +
	"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
	"`i` int(11) NOT NULL DEFAULT '0'," +
	"`s` varchar(20) NOT NULL DEFAULT ''," +
	"`s_null` varchar(20) DEFAULT NULL," +
	"`d2` varchar(20) DEFAULT NULL," +
	"`中文字段1` varchar(20) DEFAULT NULL," +
	"PRIMARY KEY (`id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

// setUp creates the database connection, recreates test tables, and inserts
// test data. It is called once before all tests via TestMain.
func setUp() {
	db, gormDB, err := mysql.NewDB()
	if err != nil {
		log.Fatalf("setUp: %v", err)
	}
	simpleDB = db

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = mysql.RecreateTable(gormDB, "diff", diffDDL)
	if err != nil {
		log.Fatalf("setUp: recreate diff: %v", err)
	}

	err = mysql.RecreateTable(gormDB, "diff_copy", diffCopyDDL)
	if err != nil {
		log.Fatalf("setUp: recreate diff_copy: %v", err)
	}

	// diff: 7 rows
	err = simpleDB.BulkInsertFromSliceSlice(ctx, "diff",
		[]string{"i", "s", "s_null", "d1"},
		[][]any{
			{0, "", nil, "nil"},
			{1, "1", "", ""},
			{2, "2", "1", "1"},
			{3, "3", "2", "2"},
			{4, "4", "3", "3"},
			{4, "4", "3", "3"},
			{4, "4", "3", "3"},
		},
		100,
	)
	if err != nil {
		log.Fatalf("setUp: insert diff: %v", err)
	}

	// diff_copy: 7 rows, slightly different data
	err = simpleDB.BulkInsertFromSliceSlice(ctx, "diff_copy",
		[]string{"i", "s", "s_null", "d2"},
		[][]any{
			{0, "", nil, nil},
			{0, "", nil, nil},
			{1, "1", "", ""},
			{2, "2", "1", "1"},
			{3, "3", "2", "2"},
			{4, "4", "3", "3"},
			{4, "4", "3", "3"},
		},
		100,
	)
	if err != nil {
		log.Fatalf("setUp: insert diff_copy: %v", err)
	}
}

// tearDown drops the test tables. It is called once after all tests via TestMain.
func tearDown() {
	if simpleDB == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = simpleDB.Drop(ctx, "diff")
	_ = simpleDB.Drop(ctx, "diff_copy")
}
