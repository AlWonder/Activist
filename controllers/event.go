package controllers

import (
	//"github.com/astaxie/beego"
	"activist_api/models"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"strconv"
)

func (c *MainController) QueryEvents() {
	var response models.QueryEventsResponse
	page, err := strconv.ParseInt(c.Input().Get("page"), 0, 64)
	if err != nil {
		c.sendErrorWithStatus("Bad page parameter", 404, 404)
	}
	response.Events, response.Count = c.getAllEvents(page)
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *MainController) QueryEventsByTag() {
	log.Println(c.Ctx.Input.Param(":tag"))
	events := c.getEventsByTag(c.Ctx.Input.Param(":tag"))
	c.Data["json"] = &events
	c.ServeJSON()
}

func (c *MainController) GetEvent() {
	var response models.GetEventResponse
	response.IsActivist = false
	response.IsJoined = false
	response.IsTimeSet = true

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	var payload jwt.MapClaims
	if payload, err = c.validateToken(); err != nil {
		log.Println("No valid token found")
	} else {
		if user := c.getUserById(int64(payload["sub"].(float64))); user != nil {
			if user.Group == 1 {
				response.IsActivist = true
				response.IsJoined = c.isJoined(user.Id, id)
			}
		}
	}

	response.Event = c.getEventById(id)
	if response.Event == nil {
		c.sendErrorWithStatus("Event not found", 404, 404)
		return
	}
	if response.Event.EventTime.IsZero() {
		response.IsTimeSet = false
	}
	response.Organizer = c.getUserById(response.Event.UserId)
	if response.Organizer == nil {
		c.sendErrorWithStatus("Organizer not found", 404, 404)
		return
	}
	response.Tags = c.getTagsByEventId(response.Event.Id)
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *MainController) QueryUserEvents() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		log.Fatal(err)
		return
	}
	events := c.getUserEvents(id)
	c.Data["json"] = &events
	c.ServeJSON()
}

func (c *MainController) QueryJoinedEvents() {
	var userId int64
	if id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64); err != nil {
		log.Fatal(err)
	} else {
		userId = id
	}

	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := c.getUserById(int64(payload["sub"].(float64)))
		if user.Group == 1 && user.Id != userId {
			c.sendErrorWithStatus("You're not allowed to do this", 403, 403)
			return
		}
	}

	events := c.getJoinedEvents(userId)
	c.Data["json"] = &events
	c.ServeJSON()
}

func (c *MainController) AddEvent() {
	var response models.AddEventResponse
	response.Ok = false
	var userId, eventId int64

	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := c.getUserById(int64(payload["sub"].(float64)))
		if user.Group == 1 {
			c.sendErrorWithStatus("You're not allowed to create events", 403, 403)
			return
		}
		userId = user.Id
	}

	// Parse the request
	var request models.AddEventRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	request.Event.UserId = userId

	// Validation
	valid := validation.Validation{}
	valid.Required(request.Event.Title, "title")
	valid.Required(request.Event.Description, "description")
	valid.Required(request.Event.EventDate, "event_date")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.appendAddEventError(&response, "Ошибка в поле "+err.Key+": "+err.Message, 400)
			log.Println("Error on " + err.Key)
		}
	}

	// Checking for having errors
	if response.Errors != nil {
		log.Println("Errors while singing up")
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// Inserting an event into the database
	if id, err := o.Insert(&request.Event); err != nil {
		c.sendError("Couldn't create an event", 14)
		return
	} else {
		eventId = id
	}
	log.Println(eventId)

	// Tags
	log.Println(request.Tags)
	tagIds := c.addTags(request.Tags)

	if ok := c.addEventTags(eventId, tagIds); !ok {
		c.sendError("Couldn't add tags to the event", 14)
		return
	}

	response.Ok = true
	response.EventId = eventId
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *MainController) AddAvatar() {
	var userId int64

	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := c.getUserById(int64(payload["sub"].(float64)))
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
		user := c.getUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	// Need checking for right event owner and correct file
	if !c.eventBelongsToUser(eventId, userId) {
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

		event := models.Event{Id: eventId}
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

func (c *MainController) EditEvent() {
	var response models.AddEventResponse
	response.Ok = false
	var userId int64

	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := c.getUserById(int64(payload["sub"].(float64)))
		if user.Group == 1 {
			c.sendErrorWithStatus("You're not allowed to edit events", 403, 403)
			return
		}
		userId = user.Id
	}

	var request models.EditEventRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	o := orm.NewOrm()

	if request.Event.UserId != userId {
		log.Println("User is not allowed to edit the event")
		c.sendErrorWithStatus("You're not allowed to edit this event", 403, 403)
		return
	}

	// Validation
	valid := validation.Validation{}
	valid.Required(request.Event.Title, "title")
	valid.Required(request.Event.Description, "description")
	valid.Required(request.Event.EventDate, "event_date")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.appendAddEventError(&response, "Ошибка в поле "+err.Key+": "+err.Message, 400)
			log.Println("Error on " + err.Key)
		}
	}

	// Checking for having errors
	if response.Errors != nil {
		log.Println("Errors while editing an event")
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	// Inserting an event into the database
	if _, err := o.Update(&request.Event); err != nil {
		log.Println(err)
		c.sendError("Couldn't edit an event", 14)
		return
	}

	// Tags
	if err := c.deleteEventTags(request.Event.Id, request.RemovedTags); err != nil {
		log.Println(err)
		c.appendAddEventError(&response, "Не удалось удалить теги.", 400)
	}

	tagIds := c.addTags(request.AddedTags)
	if ok := c.addEventTags(request.Event.Id, tagIds); !ok {
		c.appendAddEventError(&response, "Ошибка при привязке тегов", 400)
	}

	// Checking for having errors
	if response.Errors != nil {
		log.Println("Errors while editing tags")
		c.Data["json"] = response
		c.ServeJSON()
		return
	}
	response.Ok = true
	response.EventId = request.Event.Id
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *MainController) DeleteEvent() {
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
		user := c.getUserById(int64(payload["sub"].(float64)))
		if user.Group == 1 {
			c.sendErrorWithStatus("You're not allowed to delete events", 403, 403)
			return
		}
		userId = user.Id
	}

	if !c.eventBelongsToUser(eventId, userId) {
		log.Println("User is not allowed to delete the event")
		c.sendErrorWithStatus("You're not allowed to delete this event", 403, 403)
		return
	}

	o := orm.NewOrm()

	if _, err := o.Delete(&models.Event{Id: eventId}); err != nil {
		log.Println(err)
		c.sendError("Couldn't delete an event", 14)
		return
	}

	c.sendSuccess()
}

func (c *MainController) JoinEvent() {
	var userId, eventId int64
	volunteer := false
	if id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64); err != nil {
		log.Fatal(err)
	} else {
		eventId = id
	}

	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := c.getUserById(int64(payload["sub"].(float64)))
		if user.Group != 1 {
			c.sendErrorWithStatus("You're not allowed to join events", 403, 403)
			return
		}
		userId = user.Id
	}

	var request models.JoinEventRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err == nil {
		volunteer = request.AsVolunteer
	}

	var event *models.Event
	var formId int64
	var hasForm bool
	var ok bool
	// Checking if participant already has a volunteer form
	if volunteer == true {
		if event = c.getEventById(eventId); event == nil {
			c.sendErrorWithStatus("Couldn't find an event", 500, 500)
			return
		}

		if event.Volunteers == false {
			c.sendErrorWithStatus("The event doesn't provide volunteers", 403, 403)
			return
		}

		if formId, ok = c.getFormIdByOrgId(event.UserId); !ok {
			c.sendErrorWithStatus("The organizer doesn't have a volunteer form", 500, 500)
			return
		}
		hasForm = c.activistHasForm(userId, formId)
	}

	if ok = c.joinEvent(userId, eventId, volunteer); !ok {
		c.sendError("Couldn't join event", 14)
		return
	}

	// Sending successful response
	if volunteer == true {
		c.Data["json"] = models.JoinEventVolunteerResponse{Ok: true, HasForm: hasForm, OrganizerId: event.UserId}
		c.ServeJSON()
	} else {
		c.sendSuccess()
	}
}

func (c *MainController) DenyEvent() {
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
		user := c.getUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	o := orm.NewOrm()

	log.Println(userId, eventId)

	if num, err := o.Raw(`DELETE
		FROM users_events
		WHERE user_id = ? AND event_id = ?`, userId, eventId).Exec(); err != nil {
		log.Println(err)
		c.sendError("Couldn't deny an event", 14)
		return
	} else {
		log.Println(num)
	}

	c.sendSuccess()
}

func (c *MainController) getAllEvents(page int64) (*[]models.Event, int64) {
	var events []models.Event

	o := orm.NewOrm()

	if _, err := o.Raw(`SELECT events.*
					 FROM events
					 INNER JOIN users ON events.user_id=users.id
					 WHERE users.group = 2
					 ORDER BY id DESC
					 LIMIT ?, 12`,
		(page-1)*12).QueryRows(&events); err != nil {
		log.Println(err)
		return nil, 0
	}

	count, err := o.QueryTable("events").Count()
	if err != nil {
		log.Println(err)
		return nil, 0
	}
	return &events, count
}

func (c *MainController) getEventsByTag(tag string) *[]models.Event {
	var events []models.Event

	o := orm.NewOrm()

	_, err := o.Raw(`SELECT events.*
					 FROM events
					 INNER JOIN (events_tags INNER JOIN tags ON events_tags.tag_id = tags.id)
					 ON events.id = events_tags.event_id
					 WHERE tags.name = ?`,
		tag).QueryRows(&events)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &events
}

func (c *MainController) getUserEvents(userId int64) *[]models.Event {
	var events []models.Event
	o := orm.NewOrm()
	if _, err := o.Raw("SELECT * FROM events WHERE user_id = ?", userId).QueryRows(&events); err != nil {
		return nil
	}
	return &events
}

func (c *MainController) getEventById(id int64) *models.Event {
	event := models.Event{Id: id}

	o := orm.NewOrm()

	err := o.Raw("SELECT * FROM events WHERE id = ?", id).QueryRow(&event)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &event
}

func (c *MainController) getJoinedEvents(user int64) *[]models.Event {
	var events []models.Event
	o := orm.NewOrm()

	_, err := o.Raw(`SELECT e.*
					 FROM events e
					 INNER JOIN users_events ue
					 ON ue.event_id = e.id
					 WHERE ue.user_id = ?`, user).QueryRows(&events)
	if err != nil {
		log.Println("getJoinedEvents: ", err)
		return nil
	}
	return &events
}

func (c *MainController) joinEvent(user, event int64, volunteer bool) bool {
	userEvent := models.UserEvent{UserId: user, EventId: event, AsVolunteer: volunteer}

	o := orm.NewOrm()
	if _, _, err := o.ReadOrCreate(&userEvent, "UserId", "EventId", "AsVolunteer"); err != nil {
		log.Println("joinEvent: ", err)
		return false
	}
	return true
}

func (c *MainController) isJoined(user, event int64) bool {
	o := orm.NewOrm()
	userEvent := models.UserEvent{UserId: user, EventId: event}
	err := o.Read(&userEvent, "user_id", "event_id")

	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return false
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return false
	}
	return true
}

func (c *MainController) eventBelongsToUser(eventId, userId int64) bool {
	o := orm.NewOrm()
	event := models.Event{Id: eventId, UserId: userId}
	err := o.Read(&event, "id", "user_id")

	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return false
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return false
	}
	return true
}

func (c *MainController) getCover(eventId int64) string {
	o := orm.NewOrm()
	var coverPath string

	err := o.Raw(`SELECT src
		FROM images
		WHERE event_id = ?`, eventId).QueryRow(&coverPath)

	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return ""
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return ""
	}
	return coverPath
}

func (c *MainController) appendAddEventError(response *models.AddEventResponse, message string, code float64) {
	response.Errors = append(response.Errors, models.Error{
		UserMessage: message,
		Code:        code,
	})
}

/*----- I will destroy everything under this string. But later. -----*/

/*

func (c *MainController) getParticipants(eventId int64) *[]models.User {
	var users []models.User

	o := orm.NewOrm()

	_, err := o.Raw(`SELECT users.*
					 FROM users INNER JOIN (users_events INNER JOIN events ON events.id = users_events.event_id)
					 ON users.id = users_events.user_id
					 WHERE events.id = ? AND agree = 1`, eventId).QueryRows(&users)
	if err != nil {
		log.Println("getParticipants: ", err)
		return nil
	}
	log.Println(users)
	return &users
}



func (c *MainController) addTag(name string) {
	o := orm.NewOrm()
	tag := models.Tag{Name: name}
	o.Insert(&tag)
}
*/
