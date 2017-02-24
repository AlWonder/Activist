package models

import (
	"github.com/astaxie/beego/orm"
)

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

type UserEvent struct {
	Id          int64 `orm:"column(id)" json:"id"`
	UserId      int64 `orm:"column(user_id)" json:"userId"`
	EventId     int64 `orm:"column(event_id)" json:"eventId"`
	AsVolonteur bool  `orm:"column(as_volonteur);default(0)" json:"asVolonteur"`
}

type UserGroup struct {
	Id        int64  `orm:"column(id)" json:"id"`
	GroupName string `orm:"column(group_name);size(15)" json:"id"`
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
