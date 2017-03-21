package models

import (
	"encoding/json"
	"errors"
	"time"
)

type Event struct {
	Id          int64      `orm:"column(id);pk" json:"id"`
	Organizer   *User      `orm:"column(user_id);rel(fk)" json:"userId"`
	Title       string     `orm:"column(title);size(120)" json:"title,omitempty"`
	Description string     `orm:"column(description);type(text)" json:"description,omitempty"`
	CreateDate  time.Time `orm:"column(create_date);auto_now_add;type(date)" json:"createDate,omitempty"`
	EventDate   time.Time `orm:"column(event_date);type(date)" json:"eventDate,omitempty"`
	EventTime   time.Time `orm:"column(event_time);type(datetime)" json:"eventTime,omitempty"`
	Volunteers  bool       `orm:"column(volunteers);default(0)" json:"volunteers,omitempty"`
	Cover       string     `orm:"column(cover);size(128)" json:"cover,omitempty"`
}

func (e *Event) MarshalJSON() ([]byte, error) {
	type Alias Event

	loc, _ := time.LoadLocation("Etc/GMT")
	var createDate, eventDate, eventTime string
	if e.CreateDate.IsZero() {
		createDate = ""
	} else {
		createDate = e.CreateDate.Format("2006-01-02")
	}
	if e.EventDate.IsZero() {
		eventDate = ""
	} else {
		eventDate = e.EventDate.Format("2006-01-02")
	}
	if e.EventTime.IsZero() {
		eventTime = ""
	} else {
		eventTime = e.EventTime.In(loc).Format("15:04")
	}
	return json.Marshal(&struct {
		*Alias
		CreateDate string `json:"createDate,omitempty"`
		EventDate  string `json:"eventDate,omitempty"`
		EventTime  string `json:"eventTime,omitempty"`
	}{
		Alias:      (*Alias)(e),
		CreateDate: createDate,
		EventDate:  eventDate,
		EventTime:  eventTime,
	})
}

func (e *Event) UnmarshalJSON(request []byte) (err error) {
	var rawStrings map[string]interface{}

	if err := json.Unmarshal(request, &rawStrings); err != nil {
		return err
	}

	for k, v := range rawStrings {
		switch k {
		case "id":
			if id, ok := v.(float64); !ok {
				return errors.New("Bad id field")
			} else {
				e.Id = int64(id)
			}
		case "userId":
			if id, ok := v.(float64); !ok {
				return errors.New("Bad userId field")
			} else {
				e.Organizer.Id = int64(id)
			}
		case "volunteers":
			if vol, ok := v.(bool); !ok {
				return errors.New("Bad volunteers field")
			} else {
				e.Volunteers = vol
			}
		case "title":
			if title, ok := v.(string); !ok {
				return errors.New("Bad title field")
			} else {
				e.Title = title
			}
		case "description":
			if d, ok := v.(string); !ok {
				return errors.New("Bad description field")
			} else {
				e.Description = d
			}
		case "createDate":
			if createDate, err := time.Parse(time.RFC3339, v.(string)); err != nil {
				return err
			} else {
				e.CreateDate = createDate
			}
		case "eventDate":
			if eventDate, err := time.Parse("2006-01-02", v.(string)); err != nil {
				return err
			} else {
				e.EventDate = eventDate
			}
		case "eventTime":
			if timeString, ok := v.(string); ok {
				if eventTime, err := time.Parse("15:04", timeString); err != nil {
					e.EventTime = time.Time{}
				} else {
					e.EventTime = eventTime
				}
			} else {
				e.EventTime = time.Time{}
			}
		}
	}
	return
}

func (e *Event) TableName() string {
	return "events"
}
