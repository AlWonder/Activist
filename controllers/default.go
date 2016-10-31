package controllers

import (
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/orm"
	//"bee/activist/models"
	"log"
	"strconv"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.activeContent("home")
    events := c.getAllEvents(0)
    c.Data["Events"] = events
}

func (c *MainController) Profile() {
    c.activeContent("profile")
    sess := c.GetSession("activist")
    if sess != nil {
        m := sess.(map[string]interface{})
        userId := m["id"].(int64)
        log.Println(m["group"].(int64))
        if m["group"].(int64) == 1 {
            events := c.getAcceptedEvents(userId, 0)
            log.Println(events)
            c.Data["Events"] = events
        } else if m["group"].(int64) == 2 {
            events := c.getEvents(userId)
            c.Data["Events"] = events
        }
    } else {
        c.Redirect("/home", 302)
    }
}

func (c *MainController) ToHome() {
	c.Redirect("/home", 302)
}

func (c *MainController) activeContent(view string) {
    c.Layout = "basic-layout.tpl"
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["Header"] = "header.tpl"
    c.LayoutSections["Footer"] = "footer.tpl"
    c.TplName = view + ".tpl"
    sess := c.GetSession("activist")
    if sess != nil {
        c.Data["InSession"] = 1 // for login bar in header.tpl
        m := sess.(map[string]interface{})
        c.Data["Email"] = m["email"]
        c.Data["Group"] = m["group"]
        c.Data["FirstName"] = m["first_name"]
        c.Data["SecondName"] = m["second_name"]
        c.Data["LastName"] = m["last_name"]
        c.Data["Gender"] = m["gender"]
    }
}

func (c *MainController) ViewEvent() {
    c.activeContent("events/view")
    id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
    if err != nil {
        log.Fatal(err)
        return
    }
    event := c.getEvent(id)
    c.Data["Event"] = event
    sess := c.GetSession("activist")
    if sess != nil {
        m := sess.(map[string]interface{})
        c.Data["IsJoined"] = c.isJoined(m["id"].(int64), id)
    }
}

func (c *MainController) ViewParticipants() {
    c.activeContent("events/participants")
    sess := c.GetSession("activist")
    if sess != nil {
        id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
        if err != nil {
            log.Fatal(err)
            return
        }
        m := sess.(map[string]interface{})
        if c.belongsTo(id, m["id"].(int64)) {
            users := c.getParticipants(id)
            c.Data["Participants"] = users
            return
        }
    } 
    c.Redirect("/home", 302)
}


