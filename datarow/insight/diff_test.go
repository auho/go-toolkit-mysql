package insight

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	simpledb "github.com/auho/go-simple-db/v3"
	"github.com/auho/go-toolkit-mysql/datarow/insight/diff"
	"gorm.io/gorm"
)

var _simpleDB *simpledb.SimpleDB
var _gromDB *gorm.DB

func TestInsight(t *testing.T) {

	fmt.Println(fmt.Sprintf("%-20s|", "d1 [varchar]:"))
	fmt.Println(fmt.Sprintf("%-16s|", "中文字段1 [varchar]:"))

	var err error
	_simpleDB, _gromDB, err = simpledb.NewMySQLGorm("test:Test123$@tcp(127.0.0.1:3306)/test")
	if err != nil {
		t.Fatal(err)
	}

	err = _build()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_as, err := Explore(ctx, _simpleDB, "diff")
	if err != nil {
		t.Fatal(err)
	}

	for _, _s := range _as.ToStrings() {
		fmt.Println(_s)
	}

	d, err := Diff(ctx, diff.Source{
		Name: "diff",
		DB:   _simpleDB,
	}, diff.Source{
		Name: "diff_copy",
		DB:   _simpleDB,
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(strings.Join(d.DifferenceToStrings(), "\n"))
}

func _build() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := _simpleDB.Drop(ctx, "diff")
	if err != nil {
		return err
	}

	err = _simpleDB.Drop(ctx, "diff_copy")
	if err != nil {
		return err
	}

	_sql := "CREATE TABLE IF NOT EXISTS `%s` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`i` int(11) NOT NULL DEFAULT '0'," +
		"`s` varchar(20) NOT NULL DEFAULT ''," +
		"`s_null` varchar(20) DEFAULT NULL," +
		"`d1` varchar(20) DEFAULT NULL," +
		"`中文字段1` varchar(20) DEFAULT NULL," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;"

	err = _gromDB.Exec(fmt.Sprintf(_sql, "diff")).Error
	if err != nil {
		return err
	}

	_sql = "CREATE TABLE IF NOT EXISTS `%s` (" +
		"`id` int(11) unsigned NOT NULL AUTO_INCREMENT," +
		"`i` int(11) NOT NULL DEFAULT '0'," +
		"`s` varchar(20) NOT NULL DEFAULT ''," +
		"`s_null` varchar(20) DEFAULT NULL," +
		"`d2` varchar(20) DEFAULT NULL," +
		"`中文字段1` varchar(20) DEFAULT NULL," +
		"PRIMARY KEY (`id`)" +
		") ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;"

	err = _gromDB.Exec(fmt.Sprintf(_sql, "diff_copy")).Error
	if err != nil {
		return err
	}

	err = _simpleDB.BulkInsertFromSliceSlice(ctx,
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

	err = _simpleDB.BulkInsertFromSliceSlice(ctx,
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
