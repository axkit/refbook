package refbook

import "sync"

var (
	// NotFoundName returns by Name() if key not found.
	NotFoundName = "?"

	mux             sync.RWMutex
	defaultLangCode LangCode = ToLangCode("en")
)

type LangCode uint16

// Item describes JSON unmarshal destination for single language reference table.
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MultiLangItem describes JSON unmarshal destination for multi language reference table.
type MultiLangItem struct {
	ID   int               `json:"id"`
	Name map[string]string `json:"name"`
}

func ToLangCode(src string) LangCode {
	if len(src) == 0 || len(src) > 2 {
		return 0 // default language
	}
	return LangCode(uint16(src[0])<<8 | uint16(src[1]))
}

// SetDefaultLang changes default language.
func SetDefaultLang(lang string) {
	mux.Lock()
	defaultLangCode = ToLangCode(lang)
	mux.Unlock()
}
