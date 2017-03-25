package models

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/orm"
	"log"
	"time"
)

type Event struct {
	Id           int64     `orm:"column(id);pk" json:"id"`
	Organizer    *User     `orm:"column(user_id);rel(fk)" json:"userId"`
	Title        string    `orm:"column(title);size(120)" json:"title,omitempty"`
	Description  string    `orm:"column(description);type(text)" json:"description,omitempty"`
	CreateDate   time.Time `orm:"column(create_date);auto_now_add;type(date)" json:"createDate,omitempty"`
	EventDate    time.Time `orm:"column(event_date);type(date)" json:"eventDate,omitempty"`
	EventTime    time.Time `orm:"column(event_time);type(datetime)" json:"eventTime,omitempty"`
	Volunteers   bool      `orm:"column(volunteers);default(0)" json:"volunteers,omitempty"`
	Cover        string    `orm:"column(cover);size(128)" json:"cover,omitempty"`
	Participants int64     `orm:"-" json:"participants"`
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

func GetAllEvents(page int64) (*[]Event, int64) {
	var events []Event

	o := orm.NewOrm()

	if _, err := o.Raw(`SELECT events.*
					 FROM events
					 INNER JOIN users ON events.user_id=users.id
					 WHERE users.group = 2
					 ORDER BY events.id DESC
					 LIMIT ?, 12`,
		(page-1)*12).QueryRows(&events); err != nil {
		log.Println(err)
		return nil, 0
	}
	for i, _ := range events {
		if err := o.Read(events[i].Organizer); err != nil {
			log.Println(err)
			return nil, 0
		} else {
			events[i].Organizer.Email = ""
			events[i].Organizer.BirthDate = time.Time{}
		}
		if err := o.Raw(`SELECT count(*)
				FROM users_events
				WHERE event_id = ?`, events[i].Id).QueryRow(&events[i].Participants); err != nil {
			log.Println("getJoinedEvents: ", err)
		}
	}

	count, err := o.QueryTable("events").Count()
	if err != nil {
		log.Println(err)
		return nil, 0
	}
	return &events, count
}

func GetEventsByTag(tag string, page int64) (*[]Event, int64) {
	var events []Event

	o := orm.NewOrm()

	_, err := o.Raw(`SELECT events.*
					 FROM events
					 INNER JOIN (events_tags INNER JOIN tags ON events_tags.tag_id = tags.id)
					 ON events.id = events_tags.event_id
					 WHERE tags.name = ?
					 ORDER BY events.id DESC
					 LIMIT ?, 12`,
		tag, (page-1)*12).QueryRows(&events)
	if err != nil {
		log.Println(err)
		return nil, 0
	}

	for i, _ := range events {
		if err := o.Read(events[i].Organizer); err != nil {
			log.Println(err)
			return nil, 0
		} else {
			events[i].Organizer.Email = ""
			events[i].Organizer.BirthDate = time.Time{}
		}
		if err := o.Raw(`SELECT count(*)
				FROM users_events
				WHERE event_id = ?`, events[i].Id).QueryRow(&events[i].Participants); err != nil {
			log.Println("getJoinedEvents: ", err)
		}
	}

	var count int64
	err = o.Raw(`SELECT COUNT(*)
					 FROM events
					 INNER JOIN (events_tags INNER JOIN tags ON events_tags.tag_id = tags.id)
					 ON events.id = events_tags.event_id
					 WHERE tags.name = ?`, tag).QueryRow(&count)
	log.Println(count)
	if err != nil {
		log.Println(err)
		return nil, 0
	}
	return &events, count
}

func GetUserEvents(userId int64) *[]Event {
	var events []Event
	o := orm.NewOrm()
	if _, err := o.Raw("SELECT * FROM events WHERE user_id = ? ORDER BY id DESC", userId).QueryRows(&events); err != nil {
		return nil
	}
	return &events
}

func GetJoinedEvents(user int64) *[]Event {
	var events []Event
	o := orm.NewOrm()

	if _, err := o.Raw(`SELECT e.*
					 FROM events e
					 INNER JOIN users_events ue
					 ON ue.event_id = e.id
					 WHERE ue.user_id = ?`, user).QueryRows(&events); err != nil {
		log.Println("getJoinedEvents: ", err)
		return nil
	}
	return &events
}

func GetEventById(id int64) *Event {
	var event Event
	o := orm.NewOrm()

	err := o.Raw("SELECT * FROM events WHERE id = ?", id).QueryRow(&event)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &event
}

func GetSoonerEvents(limit int64) *[]Event {
	var events []Event
	o := orm.NewOrm()
	_, err := o.Raw(`SELECT *
		FROM events
		WHERE event_date >= CURDATE()
		ORDER BY event_date
		LIMIT ?`, limit).QueryRows(&events)
	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return nil
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return nil
	}
	for i, _ := range events {
		if err := o.Read(events[i].Organizer); err != nil {
			log.Println(err)
			return nil
		} else {
			events[i].Organizer.Email = ""
			events[i].Organizer.BirthDate = time.Time{}
		}
		if err := o.Raw(`SELECT count(*)
				FROM users_events
				WHERE event_id = ?`, events[i].Id).QueryRow(&events[i].Participants); err != nil {
			log.Println("getJoinedEvents: ", err)
		}
	}
	return &events
}

func GetTopFiveEventsByTags(tags *[]Tag) *[]EventsByTag {
	var eTags []EventsByTag
	o := orm.NewOrm()
	for _, tag := range *tags {
		var events []Event
		_, err := o.Raw(`SELECT e.*
			FROM events e INNER JOIN (events_tags et INNER JOIN tags t ON t.id = et.tag_id) ON et.event_id = e.id
			WHERE t.id = ?
			ORDER BY e.id DESC
			LIMIT 5`, tag.Id).QueryRows(&events)
		if err == orm.ErrNoRows {
			log.Println("No result found.")
		} else if err == orm.ErrMissPK {
			log.Println("No primary key found.")
		} else {
			for i, _ := range events {
				if err := o.Read(events[i].Organizer); err != nil {
					log.Println(err)
					return nil
				} else {
					events[i].Organizer.Email = ""
					events[i].Organizer.BirthDate = time.Time{}
				}
				if err := o.Raw(`SELECT count(*)
						FROM users_events
						WHERE event_id = ?`, events[i].Id).QueryRow(&events[i].Participants); err != nil {
					log.Println("getJoinedEvents: ", err)
				}
			}
			eTags = append(eTags, EventsByTag{Tag: tag.Name, Events: &events})
		}
	}
	return &eTags
}

func (e *Event) BelongsToUser(userId int64) bool {
	log.Println(e)
	o := orm.NewOrm()
	err := o.Raw(`SELECT *
		FROM events
		WHERE id = ? AND user_id = ?`, e.Id, userId).QueryRow(e)

	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return false
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return false
	}
	return true
}

// It would be better moved to userevent.go file. Well, later

func JoinEvent(user, event int64, volunteer bool) bool {
	userEvent := UserEvent{UserId: user, EventId: event, AsVolunteer: volunteer}

	o := orm.NewOrm()
	if _, _, err := o.ReadOrCreate(&userEvent, "UserId", "EventId", "AsVolunteer"); err != nil {
		log.Println("joinEvent: ", err)
		return false
	}
	return true
}

func DenyEvent(user, event int64) bool {
	o := orm.NewOrm()

	if _, err := o.Raw(`DELETE
		FROM users_events
		WHERE user_id = ? AND event_id = ?`, user, event).Exec(); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (e *Event) GetJoinedUsers(orgId int64) *[]JoinedUser {
	var usersEvents []UserEvent

	o := orm.NewOrm()

	if _, err := o.Raw(`SELECT *
					 FROM users_events
					 WHERE event_id = ?`, e.Id).QueryRows(&usersEvents); err != nil {
		log.Println("getJoinedUsers: ", err)
		return nil
	}

	var joinedUsers []JoinedUser

	for _, v := range usersEvents {
		user := User{Id: v.UserId}

		if err := o.Read(&user); err == orm.ErrNoRows {
			log.Println("No result found.")
		} else if err == orm.ErrMissPK {
			log.Println("No primary key found.")
		} else {
			if v.AsVolunteer {
				var formId int64
				log.Println(v.Id, orgId)
				if err := o.Raw(`SELECT fu.id
					FROM forms_users fu
					INNER JOIN (users u INNER JOIN form_templates ft ON ft.organizer_id = u.id)
					ON fu.form_id = ft.id
					WHERE fu.participant_id = ? AND u.id = ?`, v.UserId, orgId).QueryRow(&formId); err == nil {
					log.Println("Has a form")
					joinedUsers = append(joinedUsers, JoinedUser{User: user, AsVolunteer: v.AsVolunteer, FormId: formId})
				} else {
					log.Println("Doesn't have a form")
					log.Println(err)
					joinedUsers = append(joinedUsers, JoinedUser{User: user, AsVolunteer: v.AsVolunteer})
				}
			} else {
				joinedUsers = append(joinedUsers, JoinedUser{User: user, AsVolunteer: v.AsVolunteer})
			}
		}
	}
	return &joinedUsers
}
