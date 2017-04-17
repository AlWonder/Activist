package controllers

import (
	"github.com/astaxie/beego"
	"activist_api/models"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"time"
)

type EventController struct {
	beego.Controller
}

func (c *EventController) sendError(message string, code float64) {
	var response models.DefaultResponse
	response.Ok = false
	response.Error = &models.Error{ UserMessage: message, Code: code }
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *EventController) sendErrorWithStatus(message string, code float64, status int) {
	c.Ctx.Output.SetStatus(status)
	var response models.DefaultResponse
	response.Ok = false
	response.Error = &models.Error{ UserMessage: message, Code: code }
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *EventController) sendSuccess() {
	var response models.DefaultResponse
	response.Ok = true
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *EventController) QueryEvents() {
	var response models.QueryEventsResponse
	page, err := strconv.ParseInt(c.Input().Get("page"), 0, 64)
	if err != nil {
		c.sendErrorWithStatus("Bad page parameter", 404, 404)
	}
	response.Events, response.Count = models.GetAllEvents(page)
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *EventController) QueryEventsByTag() {
	var response models.QueryEventsResponse
	page, err := strconv.ParseInt(c.Input().Get("page"), 0, 64)
	if err != nil {
		c.sendErrorWithStatus("Bad page parameter", 404, 404)
	}
	log.Println(c.Ctx.Input.Param(":tag"))
	response.Events, response.Count = models.GetEventsByTag(c.Ctx.Input.Param(":tag"), page)
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *EventController) GetEvent() {
	var response models.GetEventResponse
	authenticated := false
	time.Sleep(time.Second)

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	var payload jwt.MapClaims
	if payload, err = validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println("No valid token found")
	} else {
		if user := models.GetUserById(int64(payload["sub"].(float64))); user != nil {
			authenticated = true
			if user.Group == 1 {
				response.IsJoined, response.AsVolunteer = user.IsJoined(id)
			}
		} else {
			c.sendErrorWithStatus("No user found", 403, 403)
			return
		}
	}

	if response.Event = models.GetEventById(id); response.Event == nil {
		c.sendErrorWithStatus("Event not found", 404, 404)
		return
	}
	response.Organizer = models.GetUserById(response.Event.Organizer.Id)
	if response.Organizer == nil {
		c.sendErrorWithStatus("Organizer not found", 404, 404)
		return
	}
	// Ignore email field if client isn't authenticated
	if !authenticated {
		response.Organizer.Email = ""
	}
	response.Tags = models.GetTagsByEventId(response.Event.Id)

	if response.AsVolunteer {
		// I gotta make it later, i don't know for now how many forms an organizer is allowed to have
	}

	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *EventController) QueryUserEvents() {
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		log.Fatal(err)
		return
	}
	events := models.GetUserEvents(id)
	c.Data["json"] = &events
	c.ServeJSON()
}

func (c *EventController) QueryJoinedEvents() {
	var userId int64
	if id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64); err != nil {
		log.Fatal(err)
	} else {
		userId = id
	}

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		if user.Group == 1 && user.Id != userId {
			c.sendErrorWithStatus("You're not allowed to do this", 403, 403)
			return
		}
	}

	events := models.GetJoinedEvents(userId)
	c.Data["json"] = &events
	c.ServeJSON()
}

func (c *EventController) AddEvent() {
	var response models.AddEventResponse
	response.Ok = false
	var userId, eventId int64

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
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

	organizer := models.User{Id: userId}
	request.Event.Organizer = &organizer

	// Validation
	valid := validation.Validation{}
	valid.Required(request.Event.Title, "title")
	valid.MaxSize(request.Event.Title, 120, "title")
	valid.Required(request.Event.Description, "description")
	valid.Required(request.Event.EventDate, "event_date")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			models.AppendError(&response.Errors, "Ошибка в поле "+err.Key+": "+err.Message, 400)
			log.Println("Error on " + err.Key)
		}
	}

	// Checking for having errors
	if response.Errors != nil {
		log.Println("Errors while singing up")
		c.Data["json"] = &response
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// Inserting an event into the database
	if id, err := o.Insert(&request.Event); err != nil {
		log.Println(err)
		c.sendError("Couldn't create an event", 14)
		return
	} else {
		eventId = id
	}
	log.Println(eventId)

	// Tags
	log.Println(request.Tags)
	tagIds := models.AddTags(request.Tags)

	if ok := models.AddEventTags(eventId, tagIds); !ok {
		c.sendError("Couldn't add tags to the event", 14)
		return
	}

	response.Ok = true
	response.EventId = eventId
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *EventController) EditEvent() {
	var response models.AddEventResponse
	response.Ok = false
	var userId int64

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
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

	if request.Event.Organizer.Id != userId {
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
			models.AppendError(&response.Errors, "Ошибка в поле "+err.Key+": "+err.Message, 400)
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
	if err := models.DeleteEventTags(request.Event.Id, request.RemovedTags); err != nil {
		log.Println(err)
		models.AppendError(&response.Errors, "Не удалось удалить теги.", 400)
	}

	tagIds := models.AddTags(request.AddedTags)
	if ok := models.AddEventTags(request.Event.Id, tagIds); !ok {
		models.AppendError(&response.Errors, "Ошибка при привязке тегов", 400)
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

func (c *EventController) DeleteEvent() {
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
		if user.Group == 1 {
			c.sendErrorWithStatus("You're not allowed to delete events", 403, 403)
			return
		}
		userId = user.Id
	}

	event := models.Event{Id: eventId}

	if !event.BelongsToUser(userId) {
		log.Println("User is not allowed to delete the event")
		c.sendErrorWithStatus("You're not allowed to delete this event", 403, 403)
		return
	}

	o := orm.NewOrm()

	if _, err := o.Delete(&event); err != nil {
		log.Println(err)
		c.sendError("Couldn't delete an event", 14)
		return
	}

	c.sendSuccess()
}

func (c *EventController) JoinEvent() {
	var eventId int64
	var user *models.User
	volunteer := false
	if id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64); err != nil {
		log.Fatal(err)
	} else {
		eventId = id
	}

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user = models.GetUserById(int64(payload["sub"].(float64)))
		if user.Group != 1 {
			c.sendErrorWithStatus("You're not allowed to join events", 403, 403)
			return
		}
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
		if event = models.GetEventById(eventId); event == nil {
			c.sendErrorWithStatus("Couldn't find an event", 500, 500)
			return
		}

		if event.Volunteers == false {
			c.sendErrorWithStatus("The event doesn't provide volunteers", 403, 403)
			return
		}

		if formId, ok = models.GetFormIdByOrgId(event.Organizer.Id); !ok {
			c.sendErrorWithStatus("The organizer doesn't have a volunteer form", 500, 500)
			return
		}
		hasForm = user.HasForm(formId)
	}

	if ok = models.JoinEvent(user.Id, eventId, volunteer); !ok {
		c.sendError("Couldn't join event", 14)
		return
	}

	// Sending successful response
	if volunteer == true {
		c.Data["json"] = models.JoinEventVolunteerResponse{Ok: true, HasForm: hasForm, OrganizerId: event.Organizer.Id}
		c.ServeJSON()
	} else {
		c.sendSuccess()
	}
}

func (c *EventController) DenyEvent() {
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

	if ok := models.DenyEvent(userId, eventId); !ok {
		c.sendError("Couldn't deny an event", 14)
		return
	}

	c.sendSuccess()
}
