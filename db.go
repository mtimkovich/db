package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

type StatementType int

const (
	STATEMENT_INSERT StatementType = iota
	STATEMENT_SELECT
)

var PAGE_SIZE = os.Getpagesize()

type Table struct {
	Rows      []Row
	MAX_PAGES int
}

func NewTable() *Table {
	return &Table{nil, 100}
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

type Row struct {
	Id       int
	Username string
	Email    string
}

type RawRow []byte

// Convert Row to []byte for use in pages
func (r Row) Serialize() (RawRow, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(r)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Convert RawRow back to Row. This destroys the RawRow
func (r RawRow) Deserialize() (row Row, err error) {
	buf := bytes.NewBuffer(r)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&row)

	return
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
	table := NewTable()

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
