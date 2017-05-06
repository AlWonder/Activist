package controllers

import (
	"activist_api/models"
	"crypto/rand"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/mozillazg/go-unidecode"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type FileController struct {
	beego.Controller
}

func (c *FileController) sendError(message string, code float64) {
	var response models.DefaultResponse
	response.Ok = false
	response.Error = &models.Error{UserMessage: message, Code: code}
	c.Data["json"] = &response
}

func (c *FileController) sendErrorWithStatus(message string, code float64, status int) {
	c.Ctx.Output.SetStatus(status)
	var response models.DefaultResponse
	response.Ok = false
	response.Error = &models.Error{UserMessage: message, Code: code}
	c.Data["json"] = &response
}

func (c *FileController) sendSuccess() {
	var response models.DefaultResponse
	response.Ok = true
	c.Data["json"] = &response
}

func (c *FileController) AddAvatar() {
	defer c.ServeJSON()
	var userId int64

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	log.Println("Uploading...")
	file, header, _ := c.GetFile("file") // where <<this>> is the controller and <<file>> the id of your form field
	if file != nil {
		b := make([]byte, 8)
		rand.Read(b)
		newName := fmt.Sprintf("%x", b)

		log.Println(header.Header["Content-Type"])
		if header.Header["Content-Type"][0] != "image/png" && header.Header["Content-Type"][0] != "image/jpeg" {
			c.sendError("It's not an image", 1)
			return
		}

		// save to server
		path := "static/usrfiles/user/avatar/" + newName[:2]
		_ = os.Mkdir(path, os.ModePerm)
		path += "/" + newName + ".jpg"
		err := c.SaveToFile("file", path)
		log.Println(err)

		var sFile *os.File
		if sFile, err = os.Open(path); err != nil {
			log.Println(err)
			c.sendError("Couldn't open a file", 1)
			return
		}
		if ok := transformAvatar(sFile, path); !ok {
			c.sendError("Couldn't transform an avatar", 1)
		}

		log.Println(path)

		o := orm.NewOrm()

		user := models.User{Id: userId}
		if o.Read(&user) == nil {
			user.Avatar = path[28:]
			if _, err := o.Update(&user); err == nil {
				c.sendSuccess()
				return
			}
			c.sendError("Couldn't update an avatar", 14)
			return
		}
		c.sendError("Couldn't find a user", 14)
	} else {
		c.sendError("Couldn't detect any file in the request", 1)
	}
}

func (c *FileController) AddCover() {
	defer c.ServeJSON()
	var eventId, userId int64

	eventId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	event := models.Event{Id: eventId}

	// Need checking for right event owner and correct file
	if !event.BelongsToUser(userId) {
		c.sendErrorWithStatus("You're not allowed to upload covers to this event", 403, 403)
		return
	}

	log.Println("Uploading...")
	file, header, _ := c.GetFile("file") // where <<this>> is the controller and <<file>> the id of your form field
	if file != nil {
		b := make([]byte, 8)
		rand.Read(b)
		newName := fmt.Sprintf("%x", b)

		log.Println(header.Header["Content-Type"])
		if header.Header["Content-Type"][0] != "image/png" && header.Header["Content-Type"][0] != "image/jpeg" {
			c.sendError("It's not an image", 1)
			return
		}

		// save to server
		path := "static/usrfiles/event/" + newName[:2]
		_ = os.Mkdir(path, os.ModePerm)
		path += "/" + newName[2:4]
		_ = os.Mkdir(path, os.ModePerm)
		path += "/" + newName + ".jpg"
		err := c.SaveToFile("file", path)
		log.Println(err)

		var sFile *os.File
		if sFile, err = os.Open(path); err != nil {
			log.Println(err)
			c.sendError("Couldn't open a file", 1)
			return
		}
		if ok := transformCover(sFile, path); !ok {
			c.sendError("Couldn't transform an image", 1)
		}

		log.Println(path)

		o := orm.NewOrm()

		if o.Read(&event) == nil {
			event.Cover = path[22:]
			if _, err := o.Update(&event); err == nil {
				c.sendSuccess()
				return
			}
			c.sendError("Couldn't update a cover", 14)
			return
		}
		c.sendError("Couldn't find an event", 14)
	} else {
		c.sendError("Couldn't detect a file in the request", 1)
	}
}

func (c *FileController) EditCover() {

	defer c.ServeJSON()
	var eventId, userId int64

	eventId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	event := models.GetEventById(eventId)
	if event == nil {
		c.sendErrorWithStatus("Event not found", 404, 404)
		return
	}
	log.Println(event)

	// Check for right event owner and correct file
	if !event.BelongsToUser(userId) {
		c.sendErrorWithStatus("Вам нельзя загружать обложки для этого мероприятия", 403, 403)
		return
	}

	log.Println("Uploading...")
	file, header, _ := c.GetFile("file") // where <<this>> is the controller and <<file>> the id of your form field
	if file != nil {
		b := make([]byte, 8)
		rand.Read(b)
		newName := fmt.Sprintf("%x", b)

		log.Println(header.Header["Content-Type"])
		if header.Header["Content-Type"][0] != "image/png" && header.Header["Content-Type"][0] != "image/jpeg" {
			c.sendError("It's not an image", 1)
			return
		}

		// save to server
		path := "static/usrfiles/event/" + newName[:2]
		_ = os.Mkdir(path, os.ModePerm)
		path += "/" + newName[2:4]
		_ = os.Mkdir(path, os.ModePerm)
		path += "/" + newName + ".jpg"
		err := c.SaveToFile("file", path)
		log.Println(err)

		var sFile *os.File
		if sFile, err = os.Open(path); err != nil {
			log.Println(err)
			c.sendError("Couldn't open a file", 1)
			return
		}
		if ok := transformCover(sFile, path); !ok {
			c.sendError("Couldn't transform an image", 1)
			return
		}

		log.Println(path)

		// Delete an old file
		if event.Cover != "" {
			if err := os.Remove("static/usrfiles/event/" + event.Cover); err != nil {
				log.Println(err)
			}
		}

		if ok := models.UpdateCover(event, path[22:]); !ok {
			c.sendError("Couldn't update a cover", 14)
			return
		}

		c.sendSuccess()
	} else {
		c.sendError("Couldn't detect a file in the request", 1)
	}
}

func (c *FileController) AddFormTemplate() {
	defer c.ServeJSON()
	var userId int64

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		if user.Group == 1 {
			c.sendErrorWithStatus("You're not allowed to upload templates", 403, 403)
			return
		}
		userId = user.Id
	}

	file, header, _ := c.GetFile("file") // where <<this>> is the controller and <<file>> the id of your form field
	if file != nil {

		// File extension
		ext := header.Filename[strings.LastIndex(header.Filename, "."):]

		if header.Header["Content-Type"][0] != "application/msword" &&
			header.Header["Content-Type"][0] != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" &&
			header.Header["Content-Type"][0] != "application/pdf" &&
			ext != ".doc" && ext != ".docx" && ext != ".pdf" {
			c.sendError("Unallowable file format", 1)
			log.Println("Unallowable file format")
			return
		}

		// Change a name to escape unsafe symbols
		newName := strings.Replace(url.QueryEscape(unidecode.Unidecode(header.Filename)), "+", "_", -1)
		log.Println(newName)

		// save to server
		path := "static/usrfiles/user/tpls/" + strconv.Itoa(int(userId))
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		path += "/"

		// Add an index to a file if one with the same name exists
		addName := ""
		index := 0
		exists := true
		for exists {
			if _, err = os.Stat(path + addName + newName); err == nil {
				index++
				addName = strconv.Itoa(int(index)) + "_"
			} else {
				exists = false
				err = c.SaveToFile("file", path+addName+newName)
			}
		}

		if ok := models.AddFormTemplate(userId, addName+newName); !ok {
			c.sendError("Couldn't add a template into the database", 14)
		}

		c.sendSuccess()
	} else {
		c.sendError("Couldn't detect a file in the request", 1)
	}
}

func (c *FileController) AddVolunteerForm() {
	defer c.ServeJSON()
	var userId int64

	tplId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	// Token validation
	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		// Check if the user is a participant
		if user.Group != 1 {
			c.sendErrorWithStatus("You're not allowed to upload forms", 403, 403)
			return
		}
		userId = user.Id
	}

	// Check if the form already exists
	if form := models.GetFormUser(userId, tplId); form != nil {
		c.sendErrorWithStatus("The form already exists", 409, 409)
		return
	}

	file, header, _ := c.GetFile("file")
	if file != nil {
		// File extension
		ext := header.Filename[strings.LastIndex(header.Filename, "."):]

		if header.Header["Content-Type"][0] != "application/msword" &&
			header.Header["Content-Type"][0] != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" &&
			header.Header["Content-Type"][0] != "application/pdf" &&
			ext != ".doc" && ext != ".docx" && ext != ".pdf" {
			c.sendError("Unallowable file format", 1)
			log.Println("Unallowable file format")
			return
		}

		// Change a name to escape unsafe symbols
		newName := strings.Replace(url.QueryEscape(unidecode.Unidecode(header.Filename)), "+", "_", -1)
		log.Println(newName)

		// Save to server
		path := "static/usrfiles/user/forms/" + strconv.Itoa(int(userId))
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		path += "/"

		// Add an index to a file if one with the same name exists
		addName := ""
		index := 0
		exists := true
		for exists {
			if _, err = os.Stat(path + addName + newName); err == nil {
				index++
				addName = strconv.Itoa(int(index)) + "_"
			} else {
				exists = false
				err = c.SaveToFile("file", path+addName+newName)
			}
		}

		if ok := models.AddVolunteerForm(userId, tplId, addName+newName); !ok {
			c.sendError("Couldn't add a template into the database", 14)
		}

		c.sendSuccess()
	} else {
		c.sendError("Couldn't detect a file in the request", 1)
	}
}

func (c *FileController) DeleteVolunteerForm() {
	defer c.ServeJSON()

	formId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	form := models.GetFormById(formId)
	if form == nil {
		c.sendErrorWithStatus("Form not found", 404, 404)
		return
	}

	// Token validation
	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		// Check if the user is a participant
		if user.Id != form.ParticipantId {
			c.sendErrorWithStatus("You're not allowed to delete this form", 403, 403)
			return
		}
	}

	if err := os.Remove("static/usrfiles/user/forms/" + strconv.Itoa(int(form.ParticipantId)) + "/" + form.Path); err != nil {
		log.Println(err)
		c.sendError("Couldn't delete a form", 500)
		return
	}

	if ok := models.DeleteForm(form); !ok {
		c.sendError("Couldn't delete a form", 500)
		return
	}

	c.sendSuccess()
}
