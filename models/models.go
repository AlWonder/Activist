package models

import (
    "github.com/astaxie/beego/orm"
    "time"
)

type Event struct {
	Id 				int64   	`orm:"column(id)" json:"id"` 
	UserId 			int64   	`orm:"column(user_id)" json:"user_id"`
	Name			string		`orm:"column(name);size(120)" json:"name"`
	Description		string		`orm:"column(description);type(text)" json:"description"`
	CreateDate		time.Time	`orm:"column(create_date);auto_now_add;type(date)" json:"create_date"`
	EventDate		time.Time	`orm:"column(event_date);type(date)" json:"event_date"`
	EventTime		time.Time	`orm:"column(event_time);type(datetime)" json:"event_time"`
}

func (e *Event) TableName() string {
	return "events"
}

type EventTag struct {
	Id 			int64   `orm:"column(id)" json:"id"`
	EventId 	int64   `orm:"column(event_id)" son:"event_id"`
	TagId 		int64   `orm:"column(tag_id)" json:"tag_id"`
}

func (e *EventTag) TableName() string {
	return "events_tags"
}

type FormTemplate struct {
	Id 				int64   `orm:"column(id)" json:"id"`
	OrganizerId		int64   `orm:"column(organizer_id)" json:"organizer_id"`
	TemplatePath	string	`orm:"column(template_path);size(64)" json:"template_path"`
}

func (f *FormTemplate) TableName() string {
	return "form_templates"
}

type FormUser struct {
	Id 				int64   `orm:"column(id)" json:"id"`
	ParticipantId	int64   `orm:"column(participant_id)" json:"participant_id"`
	FormId			int64   `orm:"column(form_id)" json:"form_id"`
	Path			string	`orm:"column(path);size(64)" json:"path"`
}

func (f *FormUser) TableName() string {
	return "forms_users"
}

type Image struct {
	Id 			int64   `orm:"column(id)" json:"id"`
	EventId 	int64   `orm:"column(event_id)" json:"event_id"`
	Src 		string	`orm:"column(src);size(64)" json:"src"`
}

func (i *Image) TableName() string {
	return "images"
}

type Tag struct {
	Id 				int   	`orm:"column(id);auto" json:"id"`
	Name			string	`orm:"column(name);size(30);unique" json:"name"`
}

func (t *Tag) TableName() string {
	return "tags"
}

type TagStatus struct {
	Id 		int64   `orm:"column(id)" json:"id"`
	UserId 	int64   `orm:"column(user_id)" json:"user_id"`
	TagId 	int64   `orm:"column(tag_id)" json:"tag_id"`
	Status 	bool	`orm:"column(fav_hide);default(0)" json:"fav_hide"` //0 - favorite, 1 - hidden
}

func (t *TagStatus) TableName() string {
	return "tags_status"
}

type User struct {
	Id 			int64   `orm:"column(id)" json:"id"`
	Email		string	`orm:"column(email);size(30);unique" json:"email"`
	Password	string	`orm:"column(password);size(128)" json:"password"`
	UserGroup	int64	`orm:"column(user_group);default(1)" json:"user_group"`
	FirstName	string	`orm:"column(first_name);size(25)" json:"first_name"`
	SecondName	string	`orm:"column(second_name);size(25)" json:"second_name"`
	LastName	string	`orm:"column(last_name);size(25)" json:"last_name"`
	Gender		int64	`orm:"column(gender);default(0)" json:"gender"` //0 - unknown, 1 - male, 2 - female
}

func (u *User) TableName() string {
	return "users"
}

type Who int64

const (
	prtcpnt Who = iota
	volonteur
)

type UserEvent struct {
	Id 			int64   `orm:"column(id)" json:"id"`
	UserId 		int64   `orm:"column(user_id)" json:"user_id"`
	EventId 	int64   `orm:"column(event_id)" json:"event_id"`
	Agree		bool	`orm:"column(agree);default(0)" json:"agree"`
	AsVolonteur	bool	`orm:"column(as_volonteur);default(0)" json:"as_volonteur"`
}

func (u *UserEvent) TableName() string {
	return "users_events"
}

type UserGroup struct {
	Id 			int64	`orm:"column(id)" json:"id"`
	GroupName	string	`orm:"column(group_name);size(15)" json:"id"`
}

func (u *UserGroup) TableName() string {
	return "user_groups"
}

func init() {
	orm.RegisterModel(new(Event), new(EventTag), new(FormTemplate), new(FormUser), new(Image), new(Tag), new(TagStatus), new(User), new(UserEvent), new(UserGroup))
}