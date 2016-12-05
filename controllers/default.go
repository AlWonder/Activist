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
	c.activeContent("home", "Активист")
    events := c.getAllEvents(0)
    c.Data["Events"] = events
}

func (c *MainController) Profile() {
    c.activeContent("profile", "Мой профиль")
    sess := c.GetSession("activist")
    if sess == nil {
        c.Redirect("/home", 302)
    }
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
}

func (c *MainController) ToHome() {
	c.Redirect("/home", 302)
}

func (c *MainController) activeContent(view , title string) {
    c.Layout = "basic-layout.tpl"
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["Flash"] = "flash.tpl"
    c.LayoutSections["Header"] = "header.tpl"
    c.LayoutSections["Footer"] = "footer.tpl"
    c.TplName = view + ".tpl"
    c.Data["Title"] = title;
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

// Функция для вывода только основного содержимого, без хедеров, футеров и т.д.

func (c *MainController) activeBasicContent(view string) {
    c.TplName = view + ".tpl";
}

func (c *MainController) ViewEvent() {
    id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
    if err != nil {
        log.Fatal(err)
        return
    }
    event := c.getEvent(id)
    c.Data["Event"] = event
    c.activeContent("events/view", event.Name)
    sess := c.GetSession("activist")
    if sess != nil {
        m := sess.(map[string]interface{})
        c.Data["IsJoined"] = c.isJoined(m["id"].(int64), id)
    }
}

func (c *MainController) ViewParticipants() {
    c.activeContent("events/participants", "Участники")
    sess := c.GetSession("activist")
    if sess == nil {
        c.Redirect("/home", 302)
    }
    id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
    if err != nil {
        log.Fatal(err)
        return
    }
    m := sess.(map[string]interface{})
    if c.belongsTo(id, m["id"].(int64)) {
        participants := c.getParticipants(id)
        c.Data["Participants"] = participants
    } 
}

func (c *MainController) SearchTags() {
    
    if c.Ctx.Input.Method() != "POST" {
        sess := c.GetSession("activist")
        if sess == nil {
            c.Redirect("/home", 302)
        }
        c.activeContent("searchtags", "Поиск тегов")
        return
    }
    
    c.activeBasicContent("found-tags")

    tag := c.Input().Get("tag")
    tags := c.findTags(tag)
    c.Data["Tags"] = tags

}


