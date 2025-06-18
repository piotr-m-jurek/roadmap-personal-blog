package models

type Entry struct {
	ID      string
	Title   string
	Content string
}

type Data struct {
	Page    string
	Entry   *Entry
	Entries map[string]Entry
}
