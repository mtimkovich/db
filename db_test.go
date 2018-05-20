package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var sampleRow = Row{1, "max", "max@email.com"}

func TestInsert(t *testing.T) {
	db := NewDB()

	if err := db.Insert(sampleRow); err != nil {
		t.Fatal(err)
	}

	if db.ActivePage().Rows[0] != sampleRow {
		t.Error("Inserted sampleRow not equal to Row.")
	}
}

func ExampleSelect() {
	db := NewDB()

	if err := db.Insert(sampleRow); err != nil {
		log.Fatal(err)
	}

	statement := &Statement{STATEMENT_SELECT, Row{}}
	db.Execute(statement)
	// Output: {1 max max@email.com}
}

func TestFull(t *testing.T) {
	db := NewDB()
	var err error

	for i := 0; i < os.Getpagesize()*2; i++ {
		user := fmt.Sprintf("user%v", i)
		row := Row{i, user, user + "@email.com"}

		err = db.Insert(row)

		if err != nil {
			break
		}
	}

	if err == nil {
		t.Fatal("Test should have thrown table full error.")
	}
}
