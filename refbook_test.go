package refbook

import (
	"bytes"
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestRefBook_WriteJSON(t *testing.T) {

	rb := New()
	rb.list = append(rb.list, RefBookItem{1, "A"})
	rb.list = append(rb.list, RefBookItem{2, "B"})
	rb.list = append(rb.list, RefBookItem{3, "C"})
	rb.h = calcHash(rb.list)

	buf := bytes.NewBuffer(nil)

	if _, err := rb.WriteJSON(buf); err != nil {
		t.Error(err)
	}
	if buf.String() != `{"data":[{"id":1,"name":"A"},{"id":2,"name":"B"},{"id":3,"name":"C"}],"hash":"8be"}` {
		t.Error("got wrong value", buf.String())
	}
}

func TestRefBook_Name(t *testing.T) {

	list := []RefBookItem{{1, "A"}, {2, "B"}, {3, "C"}}
	rb := New()

	rb.list = append(rb.list, list...)
	for i := range list {
		rb.idx[list[i].ID] = i
	}

	rb.h = calcHash(rb.list)

	if name := rb.Name(1); name != list[0].Name {
		t.Errorf("got %s, expected %s", name, list[0].Name)
	}

	if name := rb.Name(4); name != NotFoundName {
		t.Errorf("got %s, expected NotFoundName", name)
	}
}

func TestRefBook_LoadFromSQL(t *testing.T) {
	msg := `skipped because env var TEST_DB_CONNECTION not found.
format "user=x password=x host=127.0.0.1 port=5342 dbname=x sslmode='disable' search_path='x' bytea_output='hex'"`

	constr, ok := os.LookupEnv("TEST_DB_CONNECTION")
	if !ok {
		t.Skip(msg)
	}

	var err error
	db, err := sql.Open("postgres", constr)
	if err != nil {
		t.Error(err)
	}

	if err := db.Ping(); err != nil {
		t.Error(err)
	}

	rb := New().LoadFromSQL(db, "agents")
	if err := rb.Err(); err != nil {
		t.Error(err)
	}
}

func TestMultiLangRefBook_LoadFromSQL(t *testing.T) {
	msg := `skipped because env var TEST_DB_CONNECTION not found.
format "user=x password=x host=127.0.0.1 port=5342 dbname=x sslmode='disable' search_path='x' bytea_output='hex'"`

	constr, ok := os.LookupEnv("TEST_DB_CONNECTION")
	if !ok {
		t.Skip(msg)
	}

	var err error
	db, err := sql.Open("postgres", constr)
	if err != nil {
		t.Error(err)
	}

	if err := db.Ping(); err != nil {
		t.Error(err)
	}

	ml := NewMLRefBook().LoadFromSQL(db, "command_states")
	if err := ml.Err(); err != nil {
		t.Error(err)
	}

	t.Log(ml.Lang("ru").Name(1))
	t.Log(ml.Lang("el").Name(1))

}

// func TestRefBookWrap_WriteJSON(t *testing.T) {

// 	type Item struct {
// 		ID       int
// 		Name     map[string]string
// 		IsActive bool
// 		MaxPower float64
// 	}

// 	list := []Item{
// 		{ID: 1, Name: map[string]string{"en": "hey", "ru": "хай"}, IsActive: true, MaxPower: 100},
// 		{ID: 2, Name: map[string]string{"en": "hi", "ru": "привет"}, IsActive: true, MaxPower: 200},
// 		{ID: 3, Name: map[string]string{"ru": "привет"}, IsActive: true, MaxPower: 300},
// 	}
// 	f := func() interface{} {
// 		return list
// 	}

// 	rbw := NewRefBookerWrap(f, "Name")

// 	if _, err := rbw.WriteJSON(os.Stdout, "ru"); err != nil {
// 		t.Error(err)
// 	}
// }
func TestRefBook_LoadFromSlice(t *testing.T) {

	type Item struct {
		ID       int
		Name     string
		IsActive bool
		MaxPower float64
	}

	list := []Item{
		{ID: 1, Name: "A", IsActive: true, MaxPower: 100},
		{ID: 2, Name: "B", IsActive: true, MaxPower: 200},
		{ID: 3, Name: "C", IsActive: true, MaxPower: 300},
	}

	rb := New().LoadFromSlice(list, "ID", "Name")

	if err := rb.Err(); err != nil {
		t.Error(err)
	}

	if _, err := rb.WriteJSON(os.Stdout); err != nil {
		t.Error(err)
	}
}
func TestMultiLangRefBook_LoadFromSlice(t *testing.T) {

	type Item struct {
		ID       int
		Name     []byte
		IsActive bool
		MaxPower float64
	}

	list := []Item{
		{ID: 1, Name: []byte(`{"en": "hey", "ru": "хай"}, IsActive: true, MaxPower: 100`)},
		{ID: 2, Name: []byte(`{"en": "hi", "ru": "привет"}, IsActive: true, MaxPower: 200`)},
		{ID: 3, Name: []byte(`{"ru": "привет"}, IsActive: true, MaxPower: 300`)},
	}

	ml := NewMLRefBook().LoadFromSlice(list, "ID", "Name")

	if _, err := ml.Lang("ru").WriteJSON(os.Stdout); err != nil {
		t.Error(err)
	}
}
