package refbook

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/tidwall/gjson"
)

// FlexBook implements reference book in-memory storage.
type FlexBook struct {
	defaultLangCode LangCode
	isConcurrent    bool
	mux             sync.RWMutex
	bi              []LangCode // index in b of language hash
	book            []*Book
}

// Option holds FlexBook configuration.
type Option struct {
	lang         string
	isConcurrent bool
}

// WithDefaultLang replaces global default language.
func WithDefaultLang(lang string) func(o *Option) {
	return func(o *Option) {
		o.lang = lang
	}
}

// WithThreadSafe informs that it is changeable book.
func WithThreadSafe() func(o *Option) {
	return func(o *Option) {
		o.isConcurrent = true
	}
}

// NewFlexBook returns new reference book instance.
func NewFlexBook(f ...func(*Option)) *FlexBook {
	mux.RLock()
	lc := defaultLangCode
	mux.RUnlock()

	b := FlexBook{
		defaultLangCode: lc,
		bi:              []LangCode{lc},
		book:            []*Book{NewBook()},
	}

	o := Option{}
	for i := range f {
		f[i](&o)
	}

	if o.lang != "" {
		lc := ToLangCode(o.lang)
		b.defaultLangCode = lc
		b.bi[0] = lc
	}

	if o.isConcurrent {
		b.isConcurrent = o.isConcurrent
	}
	return &b
}

// SetThreadSafe sets flag what wraps access to internals by mutex.
func (b *FlexBook) SetThreadSafe() {
	b.isConcurrent = true
}

// Book returns pointer to the reference book associated with lang.
// Returns pointer to the book associates with default language if
// lang not found.
func (b *FlexBook) Book(lang string) *Book {
	idx := b.languageIndex(lang)
	if b.isConcurrent {
		b.mux.RLock()
	}
	res := b.book[idx]
	if b.isConcurrent {
		b.mux.RUnlock()
	}
	return res
}

func (b *FlexBook) languageIndex(lang string) int {
	lc := ToLangCode(lang)
	if lc == 0 {
		return 0
	}
	for i := range b.bi {
		if b.bi[i] == lc {
			return i
		}
	}
	return 0
}

// Name return name by id.
func (b *FlexBook) Name(lc LangCode, id int) string {
	if len(b.book) == 1 {
		res, ok := b.book[0].name(id)
		if !ok {
			return NotFoundName
		}
		return res
	}

	for i := range b.bi {
		if b.bi[i] != lc {
			continue
		}

		res, ok := b.book[i].name(id)
		if !ok && lc != 0 {
			res, ok = b.book[0].name(id)
		}
		if ok {
			return res
		}
		break
	}
	return NotFoundName
}

// IsExist returns true.
func (b *FlexBook) IsExist(id int) bool {
	if b.isConcurrent {
		b.mux.RLock()
	}
	ok := b.book[0].IsExist(id)
	if b.isConcurrent {
		b.mux.RUnlock()
	}
	return ok
}

// Len returns reference book length.
func (b *FlexBook) Len() int {
	if b.isConcurrent {
		b.mux.RLock()
	}
	res := b.book[0].Len()
	if b.isConcurrent {
		b.mux.RUnlock()
	}
	return res
}

func (b *FlexBook) BookAsJSON(lang string, dst *[]byte) {
	if b.isConcurrent {
		b.mux.RLock()
	}
	idx := b.languageIndex(lang)
	*dst = append(*dst, b.book[idx].jsonCompiled...)
	if b.isConcurrent {
		b.mux.RUnlock()
	}
}

func (b *FlexBook) Hash(lang string) uint64 {
	if b.isConcurrent {
		b.mux.RLock()
	}
	idx := b.languageIndex(lang)
	res := b.book[idx].Hash()
	if b.isConcurrent {
		b.mux.RUnlock()
	}
	return res
}

// AddItems adds items to the book.
func (b *FlexBook) AddItems(items []Item) {
	if b.isConcurrent {
		b.mux.RLock()
	}
	if len(b.book) > 1 {
		if b.isConcurrent {
			b.mux.RUnlock()
		}
		panic("AddRows called in multilang env")
	}

	db := b.book[0]
	if b.isConcurrent {
		b.mux.RUnlock()
	}

	for i := range items {
		db.Set(items[i].ID, items[i].Name)
	}
}

// AddItem adds item to the book.
func (b *FlexBook) AddItem(item Item) {
	if b.isConcurrent {
		b.mux.RLock()
	}
	if len(b.book) > 1 {
		if b.isConcurrent {
			b.mux.RUnlock()
		}
		panic("AddRows called in multilang env")
	}
	db := b.book[0]
	if b.isConcurrent {
		b.mux.RUnlock()
	}
	db.Set(item.ID, item.Name)
}

func (b *FlexBook) bookIndex(lc LangCode) int {
	for i, c := range b.bi {
		if lc == c {
			return i
		}
	}
	return -1
}

// AddMultiLangItems adds multiple items to the book.
func (b *FlexBook) AddMultiLangItems(items []MultiLangItem) {
	for i := range items {
		b.AddMultiLangItem(items[i])
	}
}

// AddMultiLangItem adds item to the book.
func (b *FlexBook) AddMultiLangItem(item MultiLangItem) {

	// make list of language codes.
	lcs := make([]LangCode, 0, len(item.Name))
	names := make(map[LangCode]string, len(lcs))

	for lang, name := range item.Name {
		if lc := ToLangCode(lang); lc > 0 {
			lcs = append(lcs, lc)
			names[lc] = name
			continue
		}
		// ignore invalid languages (lc == 0)
	}

	if b.isConcurrent {
		b.mux.Lock()
	}
	// build books in missing languages.
	for _, lc := range lcs {
		if b.bookIndex(lc) == -1 {
			b.bi = append(b.bi, lc)
			b.book = append(b.book, NewBook())
		}
	}

	for _, lc := range b.bi {
		idx := b.bookIndex(lc)
		name, ok := names[lc]
		if !ok {
			name = names[b.defaultLangCode]
			if name == "" {
				name = NotFoundName
			}
		}
		b.book[idx].Set(item.ID, name)
	}

	if b.isConcurrent {
		b.mux.Unlock()
	}
}

func (b *FlexBook) LoadFromSlice(src interface{}, attrid string, attrname string) error {

	if src == nil {
		return nil
	}

	s := reflect.ValueOf(src)
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

	nf := s.Index(0).FieldByName(attrname)
	if !nf.IsValid() {
		return fmt.Errorf("attribute %s not found", attrname)
	}

	for i := 0; i < s.Len(); i++ {
		item := s.Index(i)

		if nf.Kind() == reflect.String {
			b.AddItem(Item{ID: int(item.FieldByName(attrid).Int()), Name: item.FieldByName(attrname).String()})
		} else {
			m := map[string]string{}
			err := json.Unmarshal(item.FieldByName(attrname).Bytes(), &m)
			if err != nil {
				return err
			}
			b.AddMultiLangItem(MultiLangItem{ID: int(item.FieldByName(attrid).Int()), Name: m})
		}
	}
	return nil
}

// Parse recognizes input JSON presented as [{"id": 1, "name" : "Hello"},..] or
// [{"id":1, "name":{"en":"Hello","ru":"Привет"}},...] or mix
// [{"id":1, "name":"Hello"}, {"id":2, "name":{"en":"World", "ru":"Мир"}},...]
func (b *FlexBook) Parse(src []byte) error {

	if len(src) == 0 {
		return nil
	}

	r := gjson.GetBytes(src, "#")
	if !r.Exists() {
		return errors.New("src is not json array")
	}

	if r.Int() == 0 {
		return nil
	}

	var i int64
	ml, sl := 0, 0
	for i = 0; i < r.Int(); i++ {
		r := gjson.GetBytes(src, strconv.FormatInt(i, 10)+".name")
		switch {
		case !r.Exists():
			break
		case r.Type == gjson.JSON:
			ml++
		default:
			sl++
		}
	}
	if ml > 0 && sl > 0 {
		return errors.New("name column has different types")
	}
	if ml > 0 {

		var items []MultiLangItem
		if err := json.Unmarshal(src, &items); err != nil {
			return err
		}

		b.AddMultiLangItems(items)
		return nil
	}

	if sl > 0 {
		var items []Item
		if err := json.Unmarshal(src, &items); err != nil {
			return err
		}

		b.AddItems(items)
	}

	return nil
}
