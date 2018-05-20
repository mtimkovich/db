package main

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

var PAGE_SIZE = os.Getpagesize()

type DB struct {
	Pages     []*Page
	MAX_PAGES int
	PageIndex int
}

func NewDB() *DB {
	d := &DB{}
	d.Pages = make([]*Page, 1)
	d.Pages[0] = &Page{}
	d.MAX_PAGES = 100

	return d
}

func (d *DB) ActivePage() *Page {
	return d.Pages[d.PageIndex]
}

// Insert a row into the DB. If the page is full, make a new page! And add the row to that one.
func (d *DB) Insert(row Row) error {
	if !d.ActivePage().HasRoom(row) {
		if d.PageIndex+1 > d.MAX_PAGES {
			return errors.New("Error: Table full.")
		}

		d.PageIndex++
		d.Pages = append(d.Pages, &Page{})
	}

	d.ActivePage().Append(row)

	return nil
}

func (d *DB) Execute(statement *Statement) {
	switch statement.Type {
	case STATEMENT_INSERT:
		d.Insert(statement.RowToInsert)
	case STATEMENT_SELECT:
		for _, page := range d.Pages {
			for _, row := range page.Rows {
				fmt.Println(row)
			}
		}
	}
}

type Row struct {
	Id       int
	Username string
	Email    string
}

func (r Row) Sizeof() int {
	return int(unsafe.Sizeof(r)) + len(r.Username) + len(r.Email)
}

type Page struct {
	Rows []Row
	Size int
}

// Keep each page about the size of an OS page
func (p *Page) HasRoom(row Row) bool {
	newSize := p.Size + row.Sizeof()

	return newSize <= PAGE_SIZE
}

func (p *Page) Append(row Row) {
	p.Size += row.Sizeof()
	p.Rows = append(p.Rows, row)
}
