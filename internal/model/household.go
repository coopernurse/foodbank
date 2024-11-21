package model

import (
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

var pacLoc *time.Location

func init() {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Printf("ERROR: unable to load location America/Los_Angeles: %v", err)
	} else {
		pacLoc = loc
		fmt.Printf("Loaded location America/Los_Angeles\n")
	}
}

type Household struct {
	Id      string   `json:"id"` // Firestore document key
	Head    Person   `json:"head"`
	Members []Person `json:"members"`
}

func (h Household) Created() string {
	id, err := ulid.Parse(h.Id)
	if err == nil {
		return ulid.Time(id.Time()).In(pacLoc).Format("2006-01-02 15:04")
	}
	return ""
}
