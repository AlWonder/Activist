package controllers

import (
	"activist_api/models"
)

func (c *MainController) sendError(message string, code float64) {
	var response models.OkResponse
	response.Ok = false
	response.Errors = append(response.Errors, models.Error{ UserMessage: message, Code: code })
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *MainController) sendErrorWithStatus(message string, code float64, status int) {
	c.Ctx.Output.SetStatus(status)
	var response models.OkResponse
	response.Ok = false
	response.Errors = append(response.Errors, models.Error{ UserMessage: message, Code: code })
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *MainController) sendSuccess() {
	var response models.OkResponse
	response.Ok = true
	c.Data["json"] = response
	c.ServeJSON()
}
