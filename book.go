package refbook

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/mitchellh/hashstructure"
)

// Book implements in-memory storage of reference book in a single language.
type Book struct {
	isConcurrent bool
	mux          sync.RWMutex
	m            map[int]string
	uItems       []Item // name in upper case.
	jsonInput    struct {
		Items []Item `json:"items"`
		Hash  uint64 `json:"hash,string" hash:"ignore"`
	}
	isCompileRequired bool
	jsonCompiled      []byte
}

// NewBook returns new instance of concurrent unsafe Book.
// Use this function if you does not expect items modification
// in runtime.
func NewBook() *Book {
	return &Book{m: make(map[int]string)}
}

// NewConcurrentBook returns new instance of concurrent safe Book.
func NewConcurrentBook() *Book {
	return &Book{m: make(map[int]string), isConcurrent: true}
}

func (b *Book) name(id int) (string, bool) {
	if b.isConcurrent {
		b.mux.RLock()
	}
	res, ok := b.m[id]
	if b.isConcurrent {
		b.mux.RUnlock()
	}
	return res, ok
}

// IsExist return true if item with id exists.
func (b *Book) IsExist(id int) bool {
	if b.isConcurrent {
		b.mux.RLock()
	}
	_, ok := b.m[id]
	if b.isConcurrent {
		b.mux.RUnlock()
	}
	return ok
}

// Set inserts/update reference book item.
// Does nothing if an item exist.
func (b *Book) Set(id int, name string) {
	if b.isConcurrent {
		b.mux.Lock()
	}
	on, ok := b.m[id]
	if ok {
		if on == name {
			if b.isConcurrent {
				b.mux.Unlock()
			}
			return
		}
		b.m[id] = name
		for i := range b.jsonInput.Items {
			if b.jsonInput.Items[i].ID == id {
				b.jsonInput.Items[i].Name = name
				b.uItems[i].Name = strings.ToUpper(name)
				break
			}
		}
	} else {
		b.m[id] = name
		b.jsonInput.Items = append(b.jsonInput.Items, Item{ID: id, Name: name})
		b.uItems = append(b.uItems, Item{ID: id, Name: strings.ToUpper(name)})
	}
	b.isCompileRequired = true
	b.jsonInput.Hash = 0
	if b.isConcurrent {
		b.mux.Unlock()
	}
}

// Optimize calculates hash and pre-generates JSON.
func (b *Book) Optimize() error {
	if b.isConcurrent {
		b.mux.Lock()
		defer b.mux.Unlock()
	}
	return b.optimize()
}

func (b *Book) optimize() error {

	h, err := hashstructure.Hash(b.jsonInput.Items, &hashstructure.HashOptions{})
	if err != nil {
		return err
	}

	b.jsonInput.Hash = h
	b.jsonCompiled, err = json.Marshal(b.jsonInput)
	if err != nil {
		return err
	}

	b.isCompileRequired = false
	return nil
}

// JSON returns items as JSON array [{"id":1, "name" :"aaaa"},...].
func (b *Book) JSON() []byte {
	if b.isConcurrent {
		b.mux.RLock()
	}

	res := b.jsonCompiled

	if b.isConcurrent {
		b.mux.RUnlock()
	}
	return res
}

// Hash returns hash taken out of all items.
func (b *Book) Hash() uint64 {
	if b.isConcurrent {
		b.mux.RLock()
	}

	res := b.jsonInput.Hash

	if b.isConcurrent {
		b.mux.RUnlock()
	}
	return res
}

// Len returns items count.
func (b *Book) Len() int {
	if b.isConcurrent {
		b.mux.RLock()
	}
	res := len(b.m)
	if b.isConcurrent {
		b.mux.RUnlock()
	}
	return res
}

// Name return reference book item's name by id.
// Returns variable NotFoundName if id is not found.
func (b *Book) Name(lc LangCode, id int) string {
	if b == nil {
		return NotFoundName
	}

	if b.isConcurrent {
		b.mux.RLock()
	}

	res, ok := b.m[id]

	if b.isConcurrent {
		b.mux.RUnlock()
	}

	if ok {
		return res
	}
	return NotFoundName
}

// Traverse walks through the reference book items. Calls f() for each element.
// Aborts traverse if f() return false.
func (b *Book) Traverse(f func(id int, name string) (next bool)) {
	if b.isConcurrent {
		b.mux.RLock()
		defer b.mux.RUnlock()
	}

	for id, name := range b.m {
		if next := f(id, name); !next {
			break
		}
	}
}

// Contains adds to dst ID of reference book items if name
// contains s. The function is case unsensitive.
func (b *Book) Contains(s string, dst *[]int) {
	*dst = (*dst)[:0]
	if len(s) == 0 {
		return
	}

	us := strings.ToUpper(s)
	if b.isConcurrent {
		b.mux.RLock()
	}

	for i := range b.uItems {
		if strings.Contains(b.uItems[i].Name, us) {
			*dst = append(*dst, b.uItems[i].ID)
		}
	}

	if b.isConcurrent {
		b.mux.RUnlock()
	}
}

// LoadFromSlice init reference book with id, name pairs from any slice.
func (b *Book) LoadFromSlice(slice interface{}, attrid string, attrname string) error {

	if slice == nil {
		return nil
	}

	s := reflect.ValueOf(slice)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}

	if s.Kind() != reflect.Slice {
		return errors.New("expected argument as reference to slice")
	}

	if s.Len() == 0 {
		return nil
	}

	if !s.Index(0).FieldByName(attrid).IsValid() {
		return fmt.Errorf("attribute %s not found", attrid)
	}

	if !s.Index(0).FieldByName(attrname).IsValid() {
		return fmt.Errorf("attribute %s not found", attrname)
	}

	for i := 0; i < s.Len(); i++ {
		item := s.Index(i)
		b.Set(int(item.FieldByName(attrid).Int()),
			item.FieldByName(attrname).String())
	}

	return b.optimize()
}

// Parse parses JSON array with objects [{"id": 1, "name": "Hello"},..]
// and inits
func (b *Book) Parse(buf []byte) error {

	var items []Item

	err := json.Unmarshal(buf, &items)
	if err != nil {
		return err
	}

	if b.isConcurrent {
		b.mux.Lock()
		defer b.mux.Unlock()
	}
	for k := range b.m {
		delete(b.m, k)
	}

	b.uItems = b.uItems[:0]
	b.jsonInput.Items = b.jsonInput.Items[:0]
	b.jsonInput.Items = append(b.jsonInput.Items, items...)
	b.jsonInput.Hash = 0
	b.isCompileRequired = true
	for i := range items {
		b.uItems = append(b.uItems, Item{ID: items[i].ID, Name: strings.ToUpper(items[i].Name)})
		b.m[items[i].ID] = items[i].Name
	}

	return b.optimize()
}

func (b *Book) MarshalJSON() ([]byte, error) {
	return b.JSON(), nil
}

/*

// Add adds row to the storage.
func (ml *MultiLangBook) Add(id int, name []byte) error {
	n := map[string]string{}
	if err := json.Unmarshal(name, &n); err != nil {
		return err
	}

	for lang, v := range n {
		rb, ok := ml.b[lang]
		if !ok {
			rb = New()
			rb.Add(id, v)
			ml.b[lang] = rb
		} else {
			rb.Add(id, v)
		}
	}
	ml.h = -1
	return nil
}

func (ml *MultiLangBook) LoadFromSQL(db *sql.DB, tblname string) *MultiLangBook {

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

	for _, v := range ml.b {
		ml.h += calcHash(v.Items())
	}

	return ml
}


*/

/*
// Items returns reference book items.
func (rb *Book) Items() []Item {
	return rb.list
}

// Err return last error.
func (rb *Book) Err() error {
	return rb.err
}

// LoadFromSQL reads all rows from the table 'tblname' and keep them in the memory.
// Table has to have 'id' and 'name' columns.
func (rb *Book) LoadFromSQL(db *sql.DB, tblname string) *Book {

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
func (rb *Book) Add(id int, name string) {
	rb.list = append(rb.list, Item{ID: id, Name: name})
	rb.idx[id] = len(rb.list) - 1
	rb.h = -1
}

func calcHash(list []Item) int64 {
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
func (rb *Book) Hash() string {
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
func (rb *Book) WriteJSON(w io.Writer) (int64, error) {

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

func (ml *MultiLangBook) LoadFromSlice(slice interface{}, attrid string, attrname string) *MultiLangBook {

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
*/
