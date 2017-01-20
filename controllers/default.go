package controllers

import (
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/orm"
	//"bee/activist/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
		c.TplName = "index.html"
    c.Render()
}

/**
* Дальше идёт старый код, который по-хорошему больше не стоит использовать.
* Я оставил его потому, что пока я переписываю логику работы сервера,
* чтобы он работал только с JSON, мне нужно ориентироваться,
* как я делал это раньше
*/
/*
func (c *MainController) getSessionInfo() map[string]interface{} {
    sess := c.GetSession("activist")
    if sess == nil {
        return nil
    }
    return sess.(map[string]interface{})
}

func (c *MainController) Profile() {
    c.activeContent("profile/profile", "Мой профиль", []string{}, []string{"profile"})

    m := c.getSessionInfo()
    if m == nil {
        c.Redirect("/home", 302)
    }

    userId := m["id"].(int64)
    log.Println(m["group"].(int64))
    if m["group"].(int64) == 1 {
    } else if m["group"].(int64) == 2 {
        events := c.getEvents(userId)
        c.Data["Events"] = events
    }
}

func (c *MainController) ToHome() {
	c.Redirect("/home", 302)
}

    // Template, page title, Angular module, additional css and js
//func (c *MainController) activeContent(view, title, module string, css, js []string) {
func (c *MainController) activeContent(view, title string, css, js []string) {
    c.Layout = "basic-layout.tpl"
    c.LayoutSections = make(map[string]string)
    c.LayoutSections["Flash"] = "flash.tpl"
    c.LayoutSections["Header"] = "header.tpl"
    c.LayoutSections["Footer"] = "footer.tpl"
    c.TplName = view + ".tpl"
    c.Data["Title"] = title;
    if m := c.getSessionInfo(); m != nil {
        c.Data["InSession"] = 1 // for login bar in header.tpl
        c.Data["Email"] = m["email"]
        c.Data["Group"] = m["group"]
        c.Data["FirstName"] = m["first_name"]
        c.Data["SecondName"] = m["second_name"]
        c.Data["LastName"] = m["last_name"]
        c.Data["Gender"] = m["gender"]
    }
    if len(css) > 0 {
        c.Data["Css"] = css;
    }
    if len(js) > 0 {
        c.Data["Js"] = js;
        log.Println(js)
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
    c.activeContent("events/view", event.Name, []string{}, []string{})
    if m := c.getSessionInfo(); m != nil {
        c.Data["IsJoined"] = c.isJoined(m["id"].(int64), id)
    }
}

func (c *MainController) ViewParticipants() {
    c.activeContent("events/participants", "Участники", []string{}, []string{})
    m := c.getSessionInfo()
    if m == nil {
        c.Redirect("/home", 302)
    }
    id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
    if err != nil {
        log.Fatal(err)
        return
    }
    if c.belongsTo(id, m["id"].(int64)) {
        participants := c.getParticipants(id)
        c.Data["Participants"] = participants
    }
}

func (c *MainController) SearchTags() {

    if c.Ctx.Input.Method() != "POST" {
        if m := c.getSessionInfo(); m == nil {
            c.Redirect("/home", 302)
        }
        c.activeContent("searchtags", "Поиск тегов", []string{}, []string{})
        return
    }

    tag := c.Input().Get("tag")
    tags := c.findTags(tag)
    c.Data["json"] = &tags
    c.ServeJSON()

}


// AngularJS directives

func (c *MainController) ProfileInfo() {
    c.TplName = "profile/profile-info.tpl";

}

// JSON

func (c *MainController) AcceptedEvents() {
    m := c.getSessionInfo()
    if m == nil {
        c.Data["json"] = map[string]string{}
        c.ServeJSON()
        return
    }

    if m["group"].(int64) == 1 { // If a user is a participant
        events := c.getAcceptedEvents(m["id"].(int64), 0)
        log.Println(events)
        c.Data["json"] = &events
        c.ServeJSON()
    }
}
*/
