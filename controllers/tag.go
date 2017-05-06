package controllers

import (
	"github.com/astaxie/beego"
	"activist_api/models"
	//"encoding/json"
	//"log"
)

type TagController struct {
	beego.Controller
}

func (c *TagController) sendError(message string, code float64) {
	var response models.DefaultResponse
	response.Ok = false
	response.Error = &models.Error{ UserMessage: message, Code: code }
	c.Data["json"] = &response
}

func (c *TagController) sendErrorWithStatus(message string, code float64, status int) {
	c.Ctx.Output.SetStatus(status)
	var response models.DefaultResponse
	response.Ok = false
	response.Error = &models.Error{ UserMessage: message, Code: code }
	c.Data["json"] = &response
}

func (c *TagController) sendSuccess() {
	var response models.DefaultResponse
	response.Ok = true
	c.Data["json"] = &response
}

func (c *TagController) QueryTags() {
	defer c.ServeJSON()
	tag := c.Input().Get("query")
	tags := models.GetTags(tag)
	c.Data["json"] = &tags
}
