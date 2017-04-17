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
	TemplateId    int64  `orm:"column(template_id)" json:"templateId,omitempty"`
	Path          string `orm:"column(path);size(64)" json:"path,omitempty"`
}

func GetUserFormTemplates(userId int64) *[]FormTemplate {
	var templates []FormTemplate
	o := orm.NewOrm()
	if _, err := o.Raw("SELECT * FROM form_templates WHERE organizer_id = ?", userId).QueryRows(&templates); err != nil {
		return nil
	}
	return &templates
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

func AddFormTemplate(userId int64, path string) bool {
	formTemplate := FormTemplate{OrganizerId: userId, TemplatePath: path}

	o := orm.NewOrm()
	if _, err := o.Insert(&formTemplate); err != nil {
		log.Println("AddFormTemplate: ", err)
		return false
	}
	return true
}

func AddVolunteerForm(userId, tplId int64, path string) bool {
	form := FormUser{ParticipantId: userId, TemplateId: tplId, Path: path}

	o := orm.NewOrm()
	if _, err := o.Insert(&form); err != nil {
		log.Println("AddVolunteerForm: ", err)
		return false
	}
	return true
}

func GetFormUser(prtId, tplId int64) *FormUser {
	var form FormUser

	o := orm.NewOrm()
	err := o.Raw(`SELECT *
		FROM forms_users
		WHERE participant_id = ? AND template_id = ?`, prtId, tplId).QueryRow(&form)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &form
}

func GetFormUserById(formId int64) *FormUser {
	var form FormUser

	o := orm.NewOrm()
	err := o.Raw(`SELECT *
		FROM forms_users
		WHERE id = ?`, formId).QueryRow(&form)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &form
}

func IsAllowedToDownloadForm(userId, formId int64) bool {
	var form FormUser
	o := orm.NewOrm()
	err := o.Raw(`SELECT f.*
		FROM forms_users f INNER JOIN form_templates t ON f.template_id
		WHERE f.id = ? AND (f.participant_id = ? OR t.organizer_id = ?)`, formId, userId, userId).QueryRow(&form)
	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return false
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return false
	}
	return true
}
