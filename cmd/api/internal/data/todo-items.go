package data

import "time"

type TodoItem struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	State    string     `json:"state"`
	ClosedAt *time.Time `json:"closed_at,omitempty"`
	Tags     []string   `json:"tags"`
	Priority string     `json:"priority"`
}

/*
    ** TODO [#A] Buy Milk :home:work:
		Closed: [2023-02-18]
*/
