package tableprofile

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	simpledb "github.com/auho/go-simple-db/v3"
	"gorm.io/gorm"
)

var simpleDB *simpledb.SimpleDB
var gormDB *gorm.DB

func TestTableProfile(t *testing.T) {

	fmt.Printf("%-20s|\n", "d1 [varchar]:")
	fmt.Printf("%-16s|\n", "中文字段1 [varchar]:")

	var err error
	simpleDB, gormDB, err = simpledb.NewMySQLGorm("test:Test123$@tcp(127.0.0.1:3306)/test")
	if err != nil {
		t.Fatal(err)
	}

	err = build()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	as, err := Explore(ctx, simpleDB, "diff")
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range as.ToStrings() {
		fmt.Println(s)
	}

	d, err := CompareTables(ctx, Source{
		Name: "diff",
		DB:   simpleDB,
	}, Source{
		Name: "diff_copy",
		DB:   simpleDB,
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(strings.Join(d.Differences(), "\n"))
}

func build() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := simpleDB.Drop(ctx, "diff")
	if err != nil {
		return err
	}

	err = simpleDB.Drop(ctx, "diff_copy")
	if err != nil {
		return err
	}

	sql := "CREATE TABLE IF NOT EXISTS `%s` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`i` int(11) NOT NULL DEFAULT '0'," +
		"`s` varchar(20) NOT NULL DEFAULT ''," +
		"`s_null` varchar(20) DEFAULT NULL," +
		"`d1` varchar(20) DEFAULT NULL," +
		"`中文字段1` varchar(20) DEFAULT NULL," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;"

	err = gormDB.Exec(fmt.Sprintf(sql, "diff")).Error
	if err != nil {
		return err
	}

	sql = "CREATE TABLE IF NOT EXISTS `%s` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`i` int(11) NOT NULL DEFAULT '0'," +
		"`s` varchar(20) NOT NULL DEFAULT ''," +
		"`s_null` varchar(20) DEFAULT NULL," +
		"`d2` varchar(20) DEFAULT NULL," +
		"`中文字段1` varchar(20) DEFAULT NULL," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;"

	err = gormDB.Exec(fmt.Sprintf(sql, "diff_copy")).Error
	if err != nil {
		return err
	}

	err = simpleDB.BulkInsertFromSliceSlice(ctx,
		"diff",
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
		return err
	}

	err = simpleDB.BulkInsertFromSliceSlice(ctx,
		"diff_copy",
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
		return err
	}

	return nil
}
