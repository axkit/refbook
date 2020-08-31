package refbook

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/valyala/fasttemplate"
)

var (
	// NotFoundName returns by Name() if key not found.
	NotFoundName = "?"
	// ResponseTemplate is used as fasttemplate by method WriteJSON for reponse.
	ResponseTemplate = `{"data":[(rows)],"hash":"(hash)"}`
	// RowTemplate is used as fasttemplate by method WriteJSON for a single reference book item.
	RowTemplate = `(sep){"id":(id),"name":"(name)"}`
)

// RefBooker is an interface wraps 3 methods:
// Err
// IsExist
// Hash
type RefBooker interface {
	Err() error
	IsExist(int) bool
	Name(id int) string
	Items() []RefBookItem
	Hash() string
	Add(id int, name string)
	WriteJSON(w io.Writer) (int64, error)
}

type MultiLangRefBooker interface {
	Add(id int, name []byte) error
	Lang(lang string) RefBooker
}

type RefBookItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RefBook struct {
	err  error
	list []RefBookItem
	idx  map[int]int
	h    int64
}

type MultiLangRefBook struct {
	err error
	rb  map[string]*RefBook
	id  map[int]struct{}
	h   int64
}

var _ MultiLangRefBooker = (*MultiLangRefBook)(nil)

func NewMLRefBook() *MultiLangRefBook {
	return &MultiLangRefBook{id: make(map[int]struct{}), rb: make(map[string]*RefBook)}
}

// Lang returns RefBook corresponded to lang. If the lang is not found
// it returns nil.
func (ml *MultiLangRefBook) Lang(lang string) RefBooker {
	res := ml.rb[lang]
	return res
}

func (ml *MultiLangRefBook) IsExist(id int) bool {
	_, ok := ml.id[id]
	return ok
}

func (ml *MultiLangRefBook) Err() error {
	return ml.err
}

func (ml *MultiLangRefBook) Hash() string {
	return strconv.FormatInt(ml.h, 16)
}

// Add adds row to the storage.
func (ml *MultiLangRefBook) Add(id int, name []byte) error {
	n := map[string]string{}
	if err := json.Unmarshal(name, &n); err != nil {
		return err
	}

	for lang, v := range n {
		rb, ok := ml.rb[lang]
		if !ok {
			rb = New()
			rb.Add(id, v)
			ml.rb[lang] = rb
		} else {
			rb.Add(id, v)
		}
	}
	ml.h = -1
	return nil
}

func (ml *MultiLangRefBook) LoadFromSQL(db *sql.DB, tblname string) *MultiLangRefBook {

	qry := `select id, name from ` + tblname

	rows, err := db.Query(qry)
	if err != nil {
		ml.err = err
		return ml
	}

	defer rows.Close()

	var (
		id   int
		name []byte
	)

	for rows.Next() {
		if ml.err = rows.Scan(
			&id,
			&name,
		); ml.err != nil {
			return ml
		}

		ml.Add(id, name)
	}

	if ml.err = rows.Err(); ml.err != nil {
		return ml
	}

	for _, v := range ml.rb {
		ml.h += calcHash(v.Items())
	}

	return ml
}

// New returns new instance of single language reference book.
func New() *RefBook {
	return &RefBook{idx: make(map[int]int)}
}

// Name return reference book item's name by id.
// Returns value from variable 'NotFoundName' if id is invalid.
func (rb *RefBook) Name(id int) string {

	if rb == nil {
		return NotFoundName
	}

	idx, ok := rb.idx[id]
	if ok {
		return rb.list[idx].Name
	}
	return NotFoundName
}

// IsExist return true if reference book item is exist.
func (rb *RefBook) IsExist(id int) bool {
	_, ok := rb.idx[id]
	return ok
}

// Items returns reference book items.
func (rb *RefBook) Items() []RefBookItem {
	return rb.list
}

// Err return last error.
func (rb *RefBook) Err() error {
	return rb.err
}

// LoadFromSQL reads all rows from the table 'tblname' and keep them in the memory.
// Table has to have 'id' and 'name' columns.
func (rb *RefBook) LoadFromSQL(db *sql.DB, tblname string) *RefBook {

	qry := `select id, name from ` + tblname

	rows, err := db.Query(qry)
	if err != nil {
		rb.err = err
		return rb
	}

	defer rows.Close()

	var (
		id   int
		name *string
	)

	for rows.Next() {
		if rb.err = rows.Scan(
			&id,
			&name,
		); rb.err != nil {
			return rb
		}
		n := ""
		if name != nil {
			n = *name
		}
		rb.Add(id, n)
	}

	rb.err = rows.Err()
	if rb.err != nil {
		rb.h = calcHash(rb.list)
	}
	return rb
}

// Add adds row to the memory storage.
func (rb *RefBook) Add(id int, name string) {
	rb.list = append(rb.list, RefBookItem{ID: id, Name: name})
	rb.idx[id] = len(rb.list) - 1
	rb.h = -1
}

func calcHash(list []RefBookItem) int64 {
	var h int64
	for i := range list {
		h += int64(list[i].ID * 10)
		for _, c := range []rune(list[i].Name) {
			h += int64(c) * 11
		}
	}
	return h
}

// Hash returns hash calculed on content.
func (rb *RefBook) Hash() string {
	if rb.h == -1 {
		rb.h = calcHash(rb.list)
	}
	return strconv.FormatInt(rb.h, 16)
}

// WriteJSON writes content of reference book to io.Writer.
// Format:
// {
//  "data" : [{"ID": 1,"Name":"Open"},{"ID":2, "Name":"Close"}... ],
//  "hash" : "fgaa1"
// }
// Result can be changed using variables ResponseTemplate and ResponseRow.
func (rb *RefBook) WriteJSON(w io.Writer) (int64, error) {

	if rb == nil {
		n, err := w.Write([]byte(`{}`))
		return int64(n), err
	}

	tmpl, err := fasttemplate.NewTemplate(ResponseTemplate, "(", ")")
	if err != nil {
		return 0, err
	}

	trow, err := fasttemplate.NewTemplate(RowTemplate, "(", ")")
	if err != nil {
		return 0, err
	}

	return tmpl.ExecuteFunc(w, func(w io.Writer, tag string) (int, error) {

		if tag == "hash" {
			return w.Write([]byte(rb.Hash()))
		}

		if tag == "rows" {
			sep := ""
			var nn int

			for i := range rb.list {
				row := &rb.list[i]
				n, err := trow.ExecuteFunc(w, func(w io.Writer, tag string) (int, error) {
					if tag == "id" {
						return w.Write([]byte(strconv.FormatInt(int64(row.ID), 10)))
					}
					if tag == "name" {
						return w.Write([]byte(row.Name))
					}
					if tag == "sep" {
						return w.Write([]byte(sep))
					}
					return 0, nil
				})

				if err != nil {
					return nn, err
				}
				nn += int(n)
				sep = ","
			}
			return nn, nil
		}
		return 0, nil
	})
}

func (rb *RefBook) LoadFromSlice(slice interface{}, attrid string, attrname string) *RefBook {

	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		rb.err = errors.New("expected argument as reference to slice")
		return rb
	}

	if s.Len() == 0 {
		return rb
	}

	if !s.Index(0).FieldByName(attrid).IsValid() {
		rb.err = fmt.Errorf("attribute %s not found", attrid)
		return rb
	}

	if !s.Index(0).FieldByName(attrname).IsValid() {
		rb.err = fmt.Errorf("attribute %s not found", attrname)
		return rb
	}

	for i := 0; i < s.Len(); i++ {
		item := s.Index(i)

		rb.Add(int(item.FieldByName(attrid).Int()),
			item.FieldByName(attrname).String())
	}

	return rb
}

func (ml *MultiLangRefBook) LoadFromSlice(slice interface{}, attrid string, attrname string) *MultiLangRefBook {

	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		ml.err = errors.New("expected argument as reference to slice")
		return ml
	}

	if s.Len() == 0 {
		return ml
	}

	if !s.Index(0).FieldByName(attrid).IsValid() {
		ml.err = fmt.Errorf("attribute %s not found", attrid)
		return ml
	}

	if !s.Index(0).FieldByName(attrname).IsValid() {
		ml.err = fmt.Errorf("attribute %s not found", attrname)
		return ml
	}

	for i := 0; i < s.Len(); i++ {
		item := s.Index(i)

		ml.Add(int(item.FieldByName(attrid).Int()),
			item.FieldByName(attrname).Bytes())
	}

	return ml
}
