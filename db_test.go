package main

import (
	"fmt"
	"os"
	"testing"
)

func TestInsert(t *testing.T) {
	db := NewDB()

	row := Row{1, "max", "max@email.com"}
	err := db.Insert(row)
	if err != nil {
		t.Fatal(err)
	}

	if db.ActivePage().Rows[0] != row {
		t.Error("Inserted row not equal to row.")
	}
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
