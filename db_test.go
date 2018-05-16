package main

import (
	"testing"
)

var sampleRow = Row{1, "max", "max@max.com"}

func insertRow(table *Table, row Row) {
	insert := &Statement{STATEMENT_INSERT, row}
	table.Execute(insert)
}

func TestInsert(t *testing.T) {
	table := NewTable()
	insertRow(table, sampleRow)

	if len(table.Rows) != 1 {
		t.Errorf("Num rows should be 1, actual: %v.\n", table.Rows)
	}

	if table.Rows[0] != sampleRow {
		t.Errorf("Inserted row does not match actual row: '%v' != '%v'.\n", table.Rows, sampleRow)
	}
}

func TestStatementNegativeID(t *testing.T) {
	statement, err := NewStatement("insert -3 max max@max.com")

	if err == nil {
		t.Errorf("Creation of statement with negative ID should've failed: '%v'.\n", statement)
	}
}

func TestSerialization(t *testing.T) {
	raw, err := sampleRow.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	deserial, err := raw.Deserialize()
	if err != nil {
		t.Fatal(err)
	}

	if deserial != sampleRow {
		t.Errorf("rows not equal '%v' != '%v'.\n", sampleRow, deserial)
	}
}

func ExampleSelect() {
	table := NewTable()
	insertRow(table, sampleRow)

	table.Execute(&Statement{STATEMENT_SELECT, Row{}})
	// Output: {1 max max@max.com}
}
