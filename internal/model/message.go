package model

type Message struct {
	ID      int    `db:"id"`
	Content string `db:"content"`
}
