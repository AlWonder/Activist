package models

import (
	"github.com/astaxie/beego/orm"
)

type EventTag struct {
	Id      int64 `orm:"column(id)" json:"id"`
	Event   *Event `orm:"column(event_id);rel(fk)" json:"eventId"`
	TagId   int64 `orm:"column(tag_id)" json:"tagId,omitempty"`
}

type FormTemplate struct {
	Id           int64  `orm:"column(id)" json:"id"`
	OrganizerId  int64  `orm:"column(organizer_id)" json:"organizerId,omitempty"`
	TemplatePath string `orm:"column(template_path);size(64)" json:"templatePath,omitempty"`
}

type FormUser struct {
	Id            int64  `orm:"column(id)" json:"id"`
	ParticipantId int64  `orm:"column(participant_id)" json:"participantId,omitempty"`
	FormId        int64  `orm:"column(form_id)" json:"formId,omitempty"`
	Path          string `orm:"column(path);size(64)" json:"path,omitempty"`
}

type Image struct {
	Id      int64  `orm:"column(id)" json:"id"`
	EventId int64  `orm:"column(event_id)" json:"eventId,omitempty"`
	Src     string `orm:"column(src);size(64)" json:"src,omitempty"`
}

type Tag struct {
	Id   int64  `orm:"column(id);auto" json:"id"`
	Name string `orm:"column(name);size(30);unique" json:"name,omitempty"`
}

type TagStatus struct {
	Id     int64 `orm:"column(id)" json:"id"`
	UserId int64 `orm:"column(user_id)" json:"userId,omitempty"`
	TagId  int64 `orm:"column(tag_id)" json:"tagId,omitempty"`
	Status bool  `orm:"column(status);default(0)" json:"status,omitempty"` //0 - hidden, 1 - favorite
}

type UserEvent struct {
	Id          int64 `orm:"column(id)" json:"id"`
	UserId      int64 `orm:"column(user_id)" json:"userId,omitempty"`
	EventId     int64 `orm:"column(event_id)" json:"eventId,omitempty"`
	AsVolunteer bool  `orm:"column(as_volunteer);default(0)" json:"asVolunteer,omitempty"`
}

type UserGroup struct {
	Id        int64  `orm:"column(id)" json:"id"`
	GroupName string `orm:"column(group_name);size(15)" json:"id,omitempty"`
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
