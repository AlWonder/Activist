package models

import (
	"github.com/astaxie/beego/orm"
	"log"
)

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

func GetFormIdByOrgId(orgId int64) (int64, bool) {
	o := orm.NewOrm()
	var formId int64
	if err := o.Raw(`SELECT id
		FROM form_templates
		WHERE organizer_id = ?`, orgId).QueryRow(&formId); err != nil {
		log.Println(err)
		return 0, false
	}
	return formId, true
}
