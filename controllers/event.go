package controllers

import (
	//"github.com/astaxie/beego"
	"bee/activist/models"
	"github.com/astaxie/beego/orm"
	"log"
	//"time"
	jwt "github.com/dgrijalva/jwt-go"
	"strconv"
)

type GetEventResponse struct {
	Event      *models.Event `json:"event"`
	IsActivist bool          `json:"isActivist"`
	IsJoined   bool          `json:"isJoined"`
}

func (c *MainController) QueryEvents() {
	events := c.getAllEvents(1)
	c.Data["json"] = &events
	c.ServeJSON()
}

func (c *MainController) GetEvent() {
	var response GetEventResponse
	response.IsActivist = false
	response.IsJoined = false

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		log.Fatal(err)
		return
	}

	var payload jwt.MapClaims
	if payload, err = c.validateToken(); err != nil {
		log.Println("No valid token found")
	} else {
		if user := c.getUserById(payload["sub"].(int64)); user != nil {
			if user.Group == 1 {
				response.IsActivist = true
				response.IsJoined = c.isJoined(user.Id, id)
			}
		}
	}

	event := c.getEventById(id)
	response.Event = event
	c.Data["json"] = response
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

func (c *MainController) getAllEvents(limit int) *[]models.Event {
	var events []models.Event

	o := orm.NewOrm()

	_, err := o.Raw(`SELECT events.*
					 FROM events
					 INNER JOIN users ON events.user_id=users.id
					 WHERE users.group = 2 LIMIT ?, 10`,
		limit).QueryRows(&events)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &events
}

func (c *MainController) getUserEvents(userId int64) *[]models.Event {
	var events []models.Event

	o := orm.NewOrm()

	_, err := o.Raw("SELECT * FROM events WHERE user_id = ?", userId).QueryRows(&events)
	if err != nil {
		return nil
	}
	return &events
}

func (c *MainController) getEventById(id int64) *models.Event {
	event := models.Event{Id: id}

	o := orm.NewOrm()

	err := o.Raw("SELECT * FROM events WHERE id = ?", id).QueryRow(&event)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &event
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

/*----- I will destroy everything under this string. But later. -----*/

/*
func (c *MainController) NewEvent() {
	c.activeContent("events/new", "Новое событие", []string{}, []string{})
	flash := beego.NewFlash()

	m := c.getSessionInfo()
    if m == nil {
        c.Redirect("/home", 302)
    }

	var org int64
	org = 2
	if m["group"] != org {
		c.Redirect("/home", 302)
	}

	if c.Ctx.Input.Method() == "POST" {
		name := c.Input().Get("event-name")
		description := c.Input().Get("description")
		createDate := time.Now()
		eventDate, err := time.Parse("2006-01-02", c.Input().Get("event-date"))
		if err != nil {
			log.Println("NewEvent, eventDate: ", err)
			flash.Error("Wrong date.")
			flash.Store(&c.Controller)
			return
		}
		eventTime, err := time.Parse("15:04", c.Input().Get("event-time"))
		if err != nil {
			log.Println("NewEvent, eventTime: ", err)
		}

		log.Println("name: " + name)
		log.Println("description: " + description)
		log.Println("addDate: " + createDate.Format("2006-01-02"))
		log.Println("eventDate: " + eventDate.Format("2006-01-02"))
		log.Println("eventTime: " + eventTime.Format("2006-01-02 15:04:05"))

		valid := validation.Validation{}
		valid.MaxSize(name, 120, "name")
		valid.Required(name, "name")
		valid.Required(eventDate, "event-date")

		if valid.HasErrors() {
			errormap := []string{}
			log.Println("Validation error(s)")
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			c.Data["Errors"] = errormap
			return
		}

		event := models.Event{UserId: m["id"].(int64), Name: name, Description: description,
							  CreateDate: createDate, EventDate: eventDate, EventTime: eventTime}
		o := orm.NewOrm()


		_, err = o.Insert(&event)
		if err != nil {
			log.Println("NewEvent, data insertion: ", err)
			flash.Error("The data wasn't inserted.")
			flash.Store(&c.Controller)
			return
		}
		c.Redirect("/home", 302)
	}
}

func (c *MainController) EditEvent() {
	c.activeContent("events/edit", "Изменить событие", []string{}, []string{})
	flash := beego.NewFlash()
	sess := c.GetSession("activist")
	if sess == nil {
		c.Redirect("/home", 302)
	}
	m := sess.(map[string]interface{})
	var org int64
	org = 2
	if m["group"] != org {
		c.Redirect("/profile", 302)
	}

	eventId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
    if err != nil {
        log.Println("EditEvent: ", err)
        c.Abort("401")
    }

	o := orm.NewOrm()

	event := models.Event{Id: eventId}
	err = o.Read(&event)
	if err == orm.ErrNoRows {
	    log.Println("No result found.")
	    c.Abort("404")
	} else if err == orm.ErrMissPK {
	    log.Println("No primary key found.")
	    c.Abort("404")
	}

    log.Println(event)
    if event.UserId != m["id"].(int64) {
    	c.Abort("403")
    }

    if c.Ctx.Input.Method() == "GET" {
	    c.Data["Name"] = event.Name
		c.Data["Description"] = event.Description
		c.Data["EventDate"] = event.EventDate.Format("2006-01-02")
		if event.EventTime.IsZero() {
			c.Data["EventTime"] = ""
		} else {
			c.Data["EventTime"] = event.EventTime.Format("15:04")
		}

	} else if c.Ctx.Input.Method() == "POST" {
		name := c.Input().Get("event-name")
		description := c.Input().Get("description")
		eventDate, err := time.Parse("2006-01-02", c.Input().Get("event-date"))
		if err != nil {
			log.Println("NewEvent, eventDate: ", err)
			flash.Error("Wrong date.")
			flash.Store(&c.Controller)
			return
		}

		eventTime, err := time.Parse("15:04", c.Input().Get("event-time"))
		if err != nil {
			log.Println("NewEvent, eventTime: ", err)
		}

		log.Println("name: " + name)
		log.Println("description: " + description)
		log.Println("eventDate: " + eventDate.Format("2006-01-02"))
		log.Println("eventTime: " + eventTime.Format("2006-01-02 15:04:05"))

		valid := validation.Validation{}
		valid.MaxSize(name, 120, "name")
		valid.Required(name, "name")
		valid.Required(eventDate, "event-date")

		if valid.HasErrors() {
			errormap := []string{}
			log.Println("EditEvent: Validation error(s)")
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			c.Data["Errors"] = errormap
			return
		}

		event.Name = name
		event.Description = description
		event.EventDate = eventDate
		event.EventTime = eventTime
		if _, err := o.Update(&event); err != nil {
	        log.Println("EditEvent< data update: ", err)
	        flash.Error("The data wasn't updated.")
			flash.Store(&c.Controller)
			return
	    }
	    c.Redirect("/profile", 302)
	}
}

func (c *MainController) JoinEvent() {
	c.activeContent("events/join", "", []string{}, []string{})
	sess := c.GetSession("activist")
	if sess != nil {
		c.Redirect("/home", 302)
	}

	m := sess.(map[string]interface{})
	var prt int64
	prt = 1
	if m["group"] != prt {
		c.Redirect("/home", 302)
	}

	as, err := c.GetInt64("as")
	if err != nil {
		log.Println("JoinEvent, as: ", err)
		c.Abort("401")
	}

	if as == 1 {
		log.Println("Join as participant")
		eventId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	    if err != nil {
	        log.Println("JoinEvent, eventId: ", err)
	        c.Abort("401")
	    }

	    userId := m["id"].(int64)
	    log.Println(eventId, userId)

	    o := orm.NewOrm()

		userEvent := models.UserEvent{UserId: userId, EventId: eventId, Agree: true, AsVolonteur: false}
		_, err = o.Insert(&userEvent)
		if err != nil {
			log.Println("JoinEvent, data insertion: ", err)
			c.Abort("401")
		}

	} else if as == 2 {
		log.Println("Join as volonteur")
	}
}

func (c *MainController) DenyEvent() {
	sess := c.GetSession("activist")
	if sess == nil {
		c.Redirect("/home", 302)
	}

	eventId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
    if err != nil {
        log.Println("DenyEvent, eventId: ", err)
        c.Abort("401")
    }

	m := sess.(map[string]interface{})
	if c.isJoined(m["id"].(int64), eventId) {
		o := orm.NewOrm()
	    if num, err := o.QueryTable("users_events").Filter("user_id",
	    				 m["id"].(int64)).Filter("event_id",
	    				 eventId).Delete(); err == nil {
	        log.Println("Deleted row from users_events")
	        log.Println(num)
		} else {
			log.Println("DenyEvent, deleting: ", err)
		}
	}
	c.Redirect("/home", 302)
}

func (c *MainController) DeleteEvent() {
	sess := c.GetSession("activist")
	if sess == nil {
		c.Redirect("/home", 302)
	}

	eventId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
    if err != nil {
        log.Println("DeleteEvent: ", err)
        c.Abort("401")
    }

	m := sess.(map[string]interface{})
	if !c.belongsTo(eventId, m["id"].(int64)) {
		c.Abort("403")
	}

	o := orm.NewOrm()

	if num, err := o.Delete(&models.Event{Id: eventId}); err == nil {
		log.Println(num)
	} else {
		log.Println("DeleteEvent, deleting: ", err)
	}

	c.Redirect("/home", 302)
}

func (c *MainController) belongsTo(eventId, user int64) bool {
	o := orm.NewOrm()
	event := models.Event{Id: eventId, UserId: user}
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

func (c *MainController) getAcceptedEvents(user int64, limit int) *[]models.Event {

	var events []models.Event

	o := orm.NewOrm()

	_, err := o.Raw(`SELECT events.*
					 FROM events INNER JOIN (users_events INNER JOIN users ON users.id = users_events.user_id)
					 ON events.id = users_events.event_id
					 WHERE users.id = ? AND agree = 1
					 LIMIT ?, 10`, user, limit).QueryRows(&events)
	if err != nil {
		log.Println("getAcceptedEvents: ", err)
		return nil
	}
	log.Println(events)
	return &events
}

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
