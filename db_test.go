package main

import (
	"fmt"
	"os"
	"testing"
)

func TestInsert(t *testing.T) {
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
