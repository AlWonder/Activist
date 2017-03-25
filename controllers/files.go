package controllers

import (
	"log"
	"fmt"
	"os"
	"crypto/rand"
	//"github.com/astaxie/beego"
	"activist_api/models"
	"github.com/astaxie/beego/orm"
  "strconv"
)

func (c *MainController) AddAvatar() {
	var userId int64

	if payload, err := c.validateToken(); err != nil {
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

func (c *MainController) AddCover() {
	var eventId, userId int64

	eventId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	if payload, err := c.validateToken(); err != nil {
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
