package refbook

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

/*
func TestRefBook_WriteJSON(t *testing.T) {

	rb := New()
	rb.list = append(rb.list, Item{1, "A"})
	rb.list = append(rb.list, Item{2, "B"})
	rb.list = append(rb.list, Item{3, "C"})
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

	list := []Item{{1, "A"}, {2, "B"}, {3, "C"}}
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
*/
func TestHash(t *testing.T) {

	testCases := []struct {
		name string
		lang string
		lc   LangCode
	}{
		{"zero", "00", ToLangCode("00")},
		{"empty", "", 0},
		{"ru", "ru", ToLangCode("ru")},
	}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			lc := ToLangCode(testCases[i].lang)
			if lc != testCases[i].lc {
				t.Errorf("got %d, expected %d", lc, testCases[i].lc)
			}
		})
	}
}

func TestBook_LoadFromSlice(t *testing.T) {

	b := NewBook()

	rows := []struct {
		ID   int
		Name string
	}{
		{1, "Audi"},
		{2, "BMW"},
		{3, "Mercedes"},
	}

	if err := b.LoadFromSlice(&rows, "ID", "Name"); err != nil {
		t.Error(err)
	}

	for i := range rows {
		if !b.IsExist(rows[i].ID) {
			t.Error()
		}

		if b.Name(0, rows[i].ID) != rows[i].Name {
			t.Error()
		}

	}
}

func TestFlexLangBook_LoadFromSlice(t *testing.T) {

	b := NewFlexBook(WithThreadSafe())

	type Item struct {
		ID        int
		Name      []byte
		IsPrivate bool
		CreatedAt time.Time
		expected  map[string]string
	}

	src := []Item{
		{1, []byte(`{"en": "Hello", "ru": "Привет"}`), true, time.Now(), map[string]string{"en": "Hello", "ru": "Привет"}},
		{2, []byte(`{"en": "World"}`), true, time.Now(), map[string]string{"en": "World", "ru": "World"}},
		{3, []byte(`{"ru": "Гоу"}`), true, time.Now(), map[string]string{"en": NotFoundName, "ru": "Гоу"}},
		{4, []byte(`{}`), true, time.Now(), map[string]string{"en": NotFoundName, "ru": NotFoundName}},
	}

	if err := b.LoadFromSlice(src, "ID", "Name"); err != nil {
		t.Error(err)
	}

	t.Log(b.Len())
	t.Logf("%v\n", b.book[0].jsonInput.Items)
	t.Logf("%v\n", b.book[1].jsonInput.Items)

	for i := range src {
		if !b.IsExist(src[i].ID) {
			t.Error()
		}

		for lang, expected := range src[i].expected {
			name := b.Name(ToLangCode(lang), src[i].ID)
			if name != expected {
				t.Errorf("expected %s, got %s", expected, name)
			}
		}

	}
}

func TestBook_Contains(t *testing.T) {

	rows := []struct {
		ID   int
		Name string
	}{
		{1, "Audi"},
		{2, "BMW"},
		{3, "Mercedes"},
		{4, "Lada"},
		{5, "Peugeot"},
		{6, "Fiat"},
		{7, "Porshe"},
	}

	tc := []struct {
		substr   string
		expected []int
	}{
		{"s", []int{3, 7}},
		{"ot", []int{5}},
		{"BMW", []int{2}},
		{"xo", []int{}},
		{"", []int{}},
	}
	b := NewFlexBook(WithThreadSafe())
	if err := b.LoadFromSlice(&rows, "ID", "Name"); err != nil {
		t.Error(err)
	}

	for i := range tc {
		t.Run(tc[i].substr, func(t *testing.T) {
			var ids []int
			bk := b.Book("")
			bk.Contains(tc[i].substr, &ids)
			if len(tc[i].expected) != len(ids) {
				for j := range ids {
					if ids[j] != tc[i].expected[j] {
						t.Error()
					}
				}
			}
		})
	}

}

func TestBook_Parse(t *testing.T) {

	src := []byte(`[{"id":1,"name":"A"},{"id":2,"name":"B"},{"id":3,"name":"C"},{"id":4,"name":"D"},{"id":5,"name":"E"}, {"id":6}]`)
	tc := []Item{
		{1, "A"},
		{2, "B"},
		{3, "C"},
		{4, "D"},
		{5, "E"},
		{6, ""},
	}

	b := NewFlexBook()
	if err := b.Parse(src); err != nil {
		t.Error(err)
	}

	if b.Len() != len(tc) {
		t.Error()
	}

	for i := range tc {
		t.Run(strconv.Itoa(tc[i].ID), func(t *testing.T) {
			if name := b.Name(0, tc[i].ID); name != tc[i].Name {
				t.Errorf("expected: %s, got: %s", tc[i].Name, name)
			}
		})
	}
}

func TestFlexBook(t *testing.T) {

	mb := NewFlexBook(WithThreadSafe())

	rows := []struct {
		row      MultiLangItem
		expected map[string]string
	}{
		{MultiLangItem{1, map[string]string{"en": "Hello", "ru": "Привет"}}, map[string]string{"en": "Hello"}},
		{MultiLangItem{2, map[string]string{"en": "World"}}, map[string]string{"en": "World"}},
		{MultiLangItem{3, map[string]string{"ru": "Гоу"}}, map[string]string{"en": NotFoundName}},
		{MultiLangItem{4, map[string]string{}}, map[string]string{"en": NotFoundName, "ru": NotFoundName}},
	}

	for i := range rows {
		mb.AddMultiLangItem(rows[i].row)
	}

	fmt.Println(mb.book[0].m, mb.book[1].m)
	fmt.Println(mb.book[0].uItems, mb.book[1].uItems)

	for i := range rows {
		if !mb.IsExist(rows[i].row.ID) {
			t.Error()
		}

		for lang, expected := range rows[i].expected {
			name := mb.Name(ToLangCode(lang), rows[i].row.ID)
			if name != expected {
				t.Errorf("lang: %s, expected %s, got %s", lang, expected, name)
			}
		}
	}
}

func TestFlexBook_Parse(t *testing.T) {

	src := []byte(`[{"id":1,"name":"A"},{"id":2,"name":"B"},{"id":3,"name":"C"},{"id":4,"name":"D"},{"id":5,"name":"E"}, {"id":6}]`)
	tc := []Item{
		{1, "A"},
		{2, "B"},
		{3, "C"},
		{4, "D"},
		{5, "E"},
		{6, ""},
	}

	b := NewFlexBook()
	if err := b.Parse(src); err != nil {
		t.Error(err)
	}

	if b.Len() != len(tc) {
		t.Error()
	}

	for i := range tc {
		t.Run(strconv.Itoa(tc[i].ID), func(t *testing.T) {
			if name := b.Name(0, tc[i].ID); name != tc[i].Name {
				t.Errorf("expected: %s, got: %s", tc[i].Name, name)
			}
		})
	}
}

func TestFlexBook_ParseMultiLang(t *testing.T) {

	src := []byte(`[{"id":1,"name":{"en":"A","ru":"AA"}},
	{"id":2,"name":{"en":"B"}},{"id":3,"name":{"ru":"CC"}},{"id":4}]`)
	tc := []Item{
		{1, "A"},
		{2, "B"},
		{3, NotFoundName},
		{4, NotFoundName},
	}

	b := NewFlexBook(WithThreadSafe(), WithDefaultLang("en"))
	if err := b.Parse(src); err != nil {
		t.Error(err)
		t.Fail()
	}

	if b.Len() != len(tc) {
		t.Error()
	}

	for i := range tc {
		t.Run(strconv.Itoa(tc[i].ID), func(t *testing.T) {
			if name := b.Name(ToLangCode("en"), tc[i].ID); name != tc[i].Name {
				t.Errorf("expected: %s, got: %s", tc[i].Name, name)
			}
		})
	}
}
