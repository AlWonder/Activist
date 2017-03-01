package controllers

import (
	"log"
	//"github.com/astaxie/beego"
	"activist_api/models"
	"github.com/astaxie/beego/orm"
	"strconv"
)

func (c *MainController) GetUserInfo() {
	var response models.GetUserInfoResponse
	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.appendGetUserInfoError(&response, "Invalid token. Access denied", 401)
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = response
		c.ServeJSON()
		return
	} else {
		response.User = c.getUserById(int64(payload["sub"].(float64)))
		response.User.Password = ""
		c.Data["json"] = &response
		c.ServeJSON()
	}
}

func (c *MainController) GetJoinedUsers() {
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
		log.Println("User is not allowed to see joined users for not your events")
		c.sendErrorWithStatus("You're not allowed to see joined users for not your events", 403, 403)
		return
	}

	var response models.GetJoinedUsersResponse
	response.Ok = false

	if response.Users = c.getJoinedUsers(eventId, userId); response.Users == nil {
		c.sendError("Couldn't find joined users", 14)
		return
	}

	response.Ok = true
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *MainController) getUserByEmail(email string) *models.User {
	o := orm.NewOrm()
	user := models.User{Email: email}
	err := o.Read(&user, "email")

	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return nil
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return nil
	}
	return &user
}

func (c *MainController) getUserById(id int64) *models.User {
	o := orm.NewOrm()
	user := models.User{Id: id}
	err := o.Read(&user, "id")

	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return nil
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return nil
	}
	return &user
}

func (c *MainController) getOrgIdByEventId(eventId int64) (int64, bool) {
	o := orm.NewOrm()
	var orgId int64
	if err := o.Raw(`SELECT user_id
		FROM events
		WHERE id = ?`, eventId).QueryRow(&orgId); err != nil {
		log.Println("getOrgIdByEventId", err)
		return 0, false
	}
	return orgId, true
}

func (c *MainController) getJoinedUsers(eventId, orgId int64) *[]models.JoinedUser {
	var usersEvents []models.UserEvent

	o := orm.NewOrm()

	if _, err := o.Raw(`SELECT *
					 FROM users_events
					 WHERE event_id = ?`, eventId).QueryRows(&usersEvents); err != nil {
		log.Println("getJoinedUsers: ", err)
		return nil
	}

	var joinedUsers []models.JoinedUser

	for _, v := range usersEvents {
		user := models.User{Id: v.UserId}

		if err := o.Read(&user); err == orm.ErrNoRows {
			log.Println("No result found.")
		} else if err == orm.ErrMissPK {
			log.Println("No primary key found.")
		} else {
			if v.AsVolonteur {
				var formId int64
				log.Println(v.Id, orgId)
				if err := o.Raw(`SELECT fu.id
					FROM forms_users fu
					INNER JOIN (users u INNER JOIN form_templates ft ON ft.organizer_id = u.id)
					ON fu.form_id = ft.id
					WHERE fu.participant_id = ? AND u.id = ?`, v.UserId, orgId).QueryRow(&formId); err == nil {
					log.Println("Has a form")
					joinedUsers = append(joinedUsers, models.JoinedUser{User: user, AsVolonteur: v.AsVolonteur, FormId: formId})
				} else {
					log.Println("Doesn't have a form")
					log.Println(err)
					joinedUsers = append(joinedUsers, models.JoinedUser{User: user, AsVolonteur: v.AsVolonteur})
				}
			} else {
				joinedUsers = append(joinedUsers, models.JoinedUser{User: user, AsVolonteur: v.AsVolonteur})
			}
		}
	}
	return &joinedUsers
}

/*----- I know it's a mess below. I'll fix it. ----- */
/*
func (c *MainController) NewPassword() {
	m := c.getSessionInfo()
    if m == nil {
        c.Redirect("/home", 302)
    }

	flash := beego.NewFlash()
	c.activeContent("changepwd", "Изменить", []string{}, []string{})
	if c.Ctx.Input.Method() != "POST" {
		return
	}

	userId := m["id"].(int64)
	oldPassword := c.Input().Get("old_password")
	newPassword := c.Input().Get("new_password")
	if ok := c.changePassword(userId, oldPassword, newPassword); ok == true {
		c.Redirect("/profile", 302)
	} else {
		log.Println("Nope")
		flash.Error("Неверный старый пароль.")
		flash.Store(&c.Controller)
	}
}


func (c *MainController) changePassword(userId int64, oldPassword, newPassword string) bool {
	h := pk.HashPassword(oldPassword)

	o := orm.NewOrm()

	user := models.User{Id: userId}

	err := o.Read(&user);
	if err == orm.ErrNoRows {
    	log.Println("No result found.")
    	return false
	} else if err == orm.ErrMissPK {
	    log.Println("No primary key found.")
	    return false
	}

	var x pk.PasswordHash
	x.Hash = make([]byte, 32)
	x.Salt = make([]byte, 16)

	if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
		log.Println("ERROR:", err)
	}
	if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
		log.Println("ERROR:", err)
	}

	if !pk.MatchPassword(oldPassword, &x) {
		log.Println("Passwords don't match")
		return false
	}

	hn := pk.HashPassword(newPassword)

	user.Password = hex.EncodeToString(hn.Hash) + hex.EncodeToString(hn.Salt)
	log.Println("hex: " + hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt))

	if _, err = o.Update(&user); err != nil {
        log.Println("changePassword, data update: ", err)
		return false
    }
    return true
}
*/
