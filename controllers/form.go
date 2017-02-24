package controllers

import (
	"log"
	//"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"activist_api/models"
)

func (c *MainController) getFormIdByOrgId(orgId int64) (int64, bool) {
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

func (c *MainController) activistHasForm(prtId, formId int64) bool {
  o := orm.NewOrm()
  formUser := models.FormUser{ ParticipantId: prtId, FormId: formId }
  err := o.Read(&formUser, "participant_id", "form_id")

	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return false
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return false
	}
	return true
}
