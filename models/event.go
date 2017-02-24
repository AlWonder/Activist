package models

import (
  "time"
  "errors"
  "encoding/json"
)

type Event struct {
	Id          int64     `orm:"column(id)" json:"id"`
	UserId      int64     `orm:"column(user_id)" json:"userId"`
	Title       string    `orm:"column(title);size(120)" json:"title"`
	Description string    `orm:"column(description);type(text)" json:"description"`
	CreateDate  time.Time `orm:"column(create_date);auto_now_add;type(date)" json:"createDate"`
	EventDate   time.Time `orm:"column(event_date);type(date)" json:"eventDate"`
	EventTime   time.Time `orm:"column(event_time);type(datetime)" json:"eventTime"`
	Volonteurs  bool      `orm:"column(volonteurs);default(0)" json:"volonteurs"`
}

func (e *Event) MarshalJSON() ([]byte, error) {
	type Alias Event
	loc, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		*Alias
		EventDate string `json:"eventDate"`
		EventTime string `json:"eventTime"`
	}{
		Alias:     (*Alias)(e),
		EventDate: e.EventDate.Format("2006-01-02"),
		EventTime: e.EventTime.In(loc).Format("15:04"),
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
				e.UserId = int64(id)
			}
    case "volonteurs":
      if vol, ok := v.(bool); !ok {
        return errors.New("Bad volonteurs field")
      } else {
        e.Volonteurs = vol
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
