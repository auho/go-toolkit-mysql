package tableprofile

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestExplore(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := Explore(ctx, Source{Name: "diff", DB: simpleDB})
	if err != nil {
		t.Fatalf("Explore: %v", err)
	}

	if result.Table.Name != "diff" {
		t.Errorf("Table.Name = %q, want %q", result.Table.Name, "diff")
	}
	if result.Table.RowCount != 7 {
		t.Errorf("Table.RowCount = %d, want 7", result.Table.RowCount)
	}

	// should have 6 columns: id, i, s, s_null, d1, 中文字段1
	wantFields := []string{"id", "i", "s", "s_null", "d1", "中文字段1"}
	if len(result.FieldNames) != len(wantFields) {
		t.Fatalf("FieldNames len = %d, want %d", len(result.FieldNames), len(wantFields))
	}
	for k, fn := range wantFields {
		if result.FieldNames[k] != fn {
			t.Errorf("FieldNames[%d] = %q, want %q", k, result.FieldNames[k], fn)
		}
	}

	for _, line := range result.ToStrings() {
		fmt.Println(line)
	}
}

func TestCompareTables(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	differ, err := CompareTables(ctx,
		Source{Name: "diff", DB: simpleDB},
		Source{Name: "diff_copy", DB: simpleDB},
	)
	if err != nil {
		t.Fatalf("CompareTables: %v", err)
	}

	// tables have differences (different columns d1 vs d2, different data)
	if differ.IsOK() {
		t.Error("IsOK = true, want false for tables with differences")
	}

	diffs := differ.Differences()
	if len(diffs) == 0 {
		t.Error("Differences returned empty slice")
	}

	diffText := strings.Join(diffs, "\n")
	fmt.Println(diffText)

	// d1 exists only in left (diff), d2 exists only in right (diff_copy)
	if !strings.Contains(diffText, "d1") {
		t.Error("expected d1 in differences (column only in left)")
	}
	if !strings.Contains(diffText, "d2") {
		t.Error("expected d2 in differences (column only in right)")
	}
}
