package insight

import (
	"fmt"
	"strings"
	"testing"

	simpleDb "github.com/auho/go-simple-db/v2"
	"github.com/auho/go-toolkit/v2/mysql/datarow/insight/diff"
)

var _db *simpleDb.SimpleDB

func TestInsight(t *testing.T) {

	fmt.Println(fmt.Sprintf("%-20s|", "d1 [varchar]:"))
	fmt.Println(fmt.Sprintf("%-16s|", "中文字段1 [varchar]:"))

	var err error
	_db, err = simpleDb.NewMysql("test:Test123$@tcp(127.0.0.1:3306)/test")
	if err != nil {
		t.Fatal(err)
	}

	err = _build()
	if err != nil {
		t.Fatal(err)
	}

	_as, err := Explore(_db, "diff")
	if err != nil {
		t.Fatal(err)
	}

	for _, _s := range _as.ToStrings() {
		fmt.Println(_s)
	}

	d, err := Diff(diff.Source{
		Name: "diff",
		DB:   _db,
	}, diff.Source{
		Name: "diff_copy",
		DB:   _db,
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(strings.Join(d.DifferenceToStrings(), "\n"))
}

func _build() error {
	var err error

	err = _db.Drop("diff")
	if err != nil {
		return err
	}
	err = _db.Drop("diff_copy")
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

	err = _db.Exec(fmt.Sprintf(_sql, "diff")).Error
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

	err = _db.Exec(fmt.Sprintf(_sql, "diff_copy")).Error
	if err != nil {
		return err
	}

	err = _db.BulkInsertFromSliceSlice(
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

	err = _db.BulkInsertFromSliceSlice(
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
