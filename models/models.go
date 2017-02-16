package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
	"errors"
)

type Event struct {
	Id          int64     `orm:"column(id)" json:"id"`
	UserId      int64     `orm:"column(user_id)" json:"userId"`
	Title       string    `orm:"column(title);size(120)" json:"title"`
	Description string    `orm:"column(description);type(text)" json:"description"`
	CreateDate  time.Time `orm:"column(create_date);auto_now_add;type(date)" json:"createDate"`
	EventDate   time.Time `orm:"column(event_date);type(date)" json:"eventDate" `
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

type EventTag struct {
	Id      int64 `orm:"column(id)" json:"id"`
	EventId int64 `orm:"column(event_id)" json:"eventId"`
	TagId   int64 `orm:"column(tag_id)" json:"tagId"`
}

type FormTemplate struct {
	Id           int64  `orm:"column(id)" json:"id"`
	OrganizerId  int64  `orm:"column(organizer_id)" json:"organizerId"`
	TemplatePath string `orm:"column(template_path);size(64)" json:"templatePath"`
}

type FormUser struct {
	Id            int64  `orm:"column(id)" json:"id"`
	ParticipantId int64  `orm:"column(participant_id)" json:"participantId"`
	FormId        int64  `orm:"column(form_id)" json:"formId"`
	Path          string `orm:"column(path);size(64)" json:"path"`
}

type Image struct {
	Id      int64  `orm:"column(id)" json:"id"`
	EventId int64  `orm:"column(event_id)" json:"eventId"`
	Src     string `orm:"column(src);size(64)" json:"src"`
}

type Tag struct {
	Id   int64  `orm:"column(id);auto" json:"id"`
	Name string `orm:"column(name);size(30);unique" json:"name"`
}

type TagStatus struct {
	Id     int64 `orm:"column(id)" json:"id"`
	UserId int64 `orm:"column(user_id)" json:"userId"`
	TagId  int64 `orm:"column(tag_id)" json:"tagId"`
	Status bool  `orm:"column(fav_hide);default(0)" json:"favHide"` //0 - favorite, 1 - hidden
}

type User struct {
	Id         int64  `orm:"column(id)" json:"id"`
	Email      string `orm:"column(email);size(30);unique" json:"email"`
	Password   string `orm:"column(password);size(128)" json:"password"`
	Group      int64  `orm:"column(group);default(1)" json:"group"`
	FirstName  string `orm:"column(first_name);size(25)" json:"firstName"`
	SecondName string `orm:"column(second_name);size(25)" json:"secondName"`
	LastName   string `orm:"column(last_name);size(25)" json:"lastName"`
	Gender     int64  `orm:"column(gender);default(0)" json:"gender"` //0 - unknown, 1 - male, 2 - female
}

type UserEvent struct {
	Id          int64 `orm:"column(id)" json:"id"`
	UserId      int64 `orm:"column(user_id)" json:"userId"`
	EventId     int64 `orm:"column(event_id)" json:"eventId"`
	Agree       bool  `orm:"column(agree);default(0)" json:"agree"`
	AsVolonteur bool  `orm:"column(as_volonteur);default(0)" json:"asVolonteur"`
}

type UserGroup struct {
	Id        int64  `orm:"column(id)" json:"id"`
	GroupName string `orm:"column(group_name);size(15)" json:"id"`
}

func (e *Event) TableName() string {
	return "events"
}

func (e *EventTag) TableName() string {
	return "events_tags"
}

func (f *FormTemplate) TableName() string {
	return "form_templates"
}

func (f *FormUser) TableName() string {
	return "forms_users"
}

func (i *Image) TableName() string {
	return "images"
}

func (t *Tag) TableName() string {
	return "tags"
}

func (t *TagStatus) TableName() string {
	return "tags_status"
}

func (u *User) TableName() string {
	return "users"
}

func (u *UserEvent) TableName() string {
	return "users_events"
}

func (u *UserGroup) TableName() string {
	return "user_groups"
}

func init() {
	orm.RegisterModel(new(Event), new(EventTag), new(FormTemplate),
		new(FormUser), new(Image), new(Tag),
		new(TagStatus), new(User), new(UserEvent),
		new(UserGroup))
}
