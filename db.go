package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type StatementType int

const (
	STATEMENT_INSERT StatementType = iota
	STATEMENT_SELECT
)

type Table struct {
	Rows []Row
}

type Row struct {
	Id       int
	Username string
	Email    string
}

type Statement struct {
	Type        StatementType
	RowToInsert Row
}

// Parses input into Statement object
func NewStatement(input string) (*Statement, error) {
	s := &Statement{}

	if strings.HasPrefix(input, "select") {
		s.Type = STATEMENT_SELECT
	} else if strings.HasPrefix(input, "insert") {
		s.Type = STATEMENT_INSERT
		numArgs, err := fmt.Sscanf(input, "insert %d %s %s",
			&s.RowToInsert.Id,
			&s.RowToInsert.Username,
			&s.RowToInsert.Email)

		if numArgs < 3 || err != nil {
			return nil, fmt.Errorf("Syntax error: %v", err)
		} else if s.RowToInsert.Id < 1 {
			return nil, fmt.Errorf("ID must be positive.")
		}
	} else {
		return nil, fmt.Errorf("Unrecognized keyword at start of '%v'.", input)
	}

	return s, nil
}

func (t *Table) Execute(statement *Statement) {
	switch statement.Type {
	case STATEMENT_INSERT:
		t.Rows = append(t.Rows, statement.RowToInsert)
	case STATEMENT_SELECT:
		for _, row := range t.Rows {
			fmt.Println(row)
		}
	}
}

func prompt(ps string) string {
	fmt.Print(ps)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func doMetaCommand(input string) error {
	switch input {
	case ".exit":
		os.Exit(0)
	default:
		return fmt.Errorf("Unrecognized command '%v'.", input)
	}

	return nil
}

func main() {
	table := &Table{}

	for {
		input := prompt("db> ")

		// Metacommands
		if input[0] == '.' {
			if err := doMetaCommand(input); err != nil {
				fmt.Println(err)
			}

			continue
		}

		statement, err := NewStatement(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		table.Execute(statement)
		fmt.Println("Executed.")
	}
}
