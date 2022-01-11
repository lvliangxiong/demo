package model

type Message struct {
	Id   int    `json:"id,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Cost int    `json:"cost,omitempty"`
}
