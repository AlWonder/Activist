package controllers

import (
  "github.com/astaxie/beego"
	"log"
  "strconv"
	//"os"
	"activist_api/models"
)

type FormController struct {
	beego.Controller
}

func (c *FormController) sendError(message string, code float64) {
	var response models.DefaultResponse
	response.Ok = false
	response.Error = &models.Error{ UserMessage: message, Code: code }
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *FormController) sendErrorWithStatus(message string, code float64, status int) {
	c.Ctx.Output.SetStatus(status)
	var response models.DefaultResponse
	response.Ok = false
	response.Error = &models.Error{ UserMessage: message, Code: code }
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *FormController) sendSuccess() {
	var response models.DefaultResponse
	response.Ok = true
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *FormController) QueryUserFormTemplates() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		log.Fatal(err)
		return
	}

  if _, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	}

	templates := models.GetUserFormTemplates(id)
	c.Data["json"] = &templates
	c.ServeJSON()
}

func (c *FormController) QueryUserForms() {
  var userId int64

  if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
    user := models.GetUserById(int64(payload["sub"].(float64)))
		if user.Group != 1 {
			c.sendErrorWithStatus("You're not a participant", 403, 403)
			return
		}
    userId = user.Id
  }

	forms := models.GetUserForms(userId)
	c.Data["json"] = &forms
	c.ServeJSON()
}
