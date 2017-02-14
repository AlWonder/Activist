package controllers

import (
	"log"
	//"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"activist_api/models"
)

func (c *MainController) GetUserInfo() {
	var response GetUserInfoResponse
	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.appendGetUserInfoError(&response, "Invalid token. Access denied", 401)
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = response
		c.ServeJSON()
		return
	} else {
		user := c.getUserById(int64(payload["sub"].(float64)))
		log.Println(user)
		response.User = UserInfo {
                  Email:       &user.Email,
                  Group:       &user.Group,
                  FirstName:   &user.FirstName,
                  SecondName:  &user.SecondName,
                  LastName:    &user.LastName,
                  Gender:      &user.Gender,
		}
		c.Data["json"] = response
		c.ServeJSON()
	}
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
