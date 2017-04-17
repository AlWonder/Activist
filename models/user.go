package models

import (
  "time"
  "errors"
  "encoding/json"
	"github.com/astaxie/beego/orm"
	"log"
)

type User struct {
	Id         int64     `orm:"column(id);pk" json:"id"`
	Email      string    `orm:"column(email);size(30);unique" json:"email,omitempty"`
	Password   string    `orm:"column(password);size(128)" json:"-"`
	Group      int64     `orm:"column(group);default(1)" json:"group,omitempty"`
	FirstName  string    `orm:"column(first_name);size(25)" json:"firstName,omitempty"`
	SecondName string    `orm:"column(second_name);size(25)" json:"secondName,omitempty"`
	LastName   string    `orm:"column(last_name);size(25)" json:"lastName,omitempty"`
	BirthDate  time.Time `orm:"column(birth_date);type(date)" json:"birthDate,omitempty"`
	Gender     int64     `orm:"column(gender);default(0)" json:"gender,omitempty"` //0 - unknown, 1 - male, 2 - female
  Avatar     string    `orm:"column(avatar);size(30)" json:"avatar,omitempty"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User

  var birthDate string
  if u.BirthDate.IsZero() {
		birthDate = ""
	} else {
		birthDate = u.BirthDate.Format("2006-01-02")
	}
	return json.Marshal(&struct {
		*Alias
		BirthDate string `json:"birthDate,omitempty"`
	}{
		Alias:     (*Alias)(u),
		BirthDate: birthDate,
	})
}

func (u *User) UnmarshalJSON(request []byte) (err error) {
	var rawStrings map[string]interface{}

	if err := json.Unmarshal(request, &rawStrings); err != nil {
		return err
	}

	for k, v := range rawStrings {
		switch k {
		case "id":
			if id, ok := v.(float64); !ok {
				return errors.New("Bad id field")
			} else {
				u.Id = int64(id)
			}
		case "email":
			if email, ok := v.(string); !ok {
				return errors.New("Bad email field")
			} else {
				u.Email = email
			}
    case "password":
      if password, ok := v.(string); !ok {
        return errors.New("Bad password field")
      } else {
        u.Password = password
      }
    case "group":
      if group, ok := v.(float64); !ok {
        return errors.New("Bad group field")
      } else {
        u.Group = int64(group)
      }
    case "firstName":
      if fName, ok := v.(string); !ok {
        return errors.New("Bad firstName field")
      } else {
        u.FirstName = fName
      }
    case "secondName":
      if sName, ok := v.(string); !ok {
        return errors.New("Bad secondName field")
      } else {
        u.SecondName = sName
      }
    case "lastName":
      if lName, ok := v.(string); !ok {
        return errors.New("Bad lastName field")
      } else {
        u.LastName = lName
      }
		case "birthDate":
			if birthDate, err := time.Parse("2006-01-02", v.(string)); err != nil {
				return err
			} else {
				u.BirthDate = birthDate
			}
    case "gender":
      if gender, ok := v.(float64); !ok {
        return errors.New("Bad gender field")
      } else {
        u.Gender = int64(gender)
      }
		}
	}
	return
}

func GetUserByEmail(email string) *User {
	o := orm.NewOrm()
	user := User{Email: email}
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

func GetUserById(id int64) *User {
	o := orm.NewOrm()
	user := User{Id: id}
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

func (u *User) IsJoined(event int64) (bool, bool) {
	o := orm.NewOrm()
	userEvent := UserEvent{UserId: u.Id, EventId: event}
	err := o.Read(&userEvent, "user_id", "event_id")

	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return false, false
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return false, false
	}
	return true, userEvent.AsVolunteer
}

func (u *User) HasForm(tplId int64) bool {
  o := orm.NewOrm()
  formUser := FormUser{ ParticipantId: u.Id, TemplateId: tplId }
  err := o.Read(&formUser, "participant_id", "form_id")

	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return false
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return false
	}
	return true
}
