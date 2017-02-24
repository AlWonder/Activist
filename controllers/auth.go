package controllers

import (
	"activist_api/models"
	pk "activist_api/utilities/pbkdf2"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"time"
)

var privateKey = []byte("pisos")

func (c *MainController) Login() {
	var request models.LoginRequest
	var response models.LoginResponse

	// Checking for a correct JSON request. If not, throw an error to a client
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	valid := validation.Validation{}
	valid.Required(request.Email, "email")
	valid.Required(request.Password, "password")
	valid.Email(request.Email, "email")
	valid.MaxSize(request.Email, 30, "email")
	valid.MaxSize(request.Password, 30, "email")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.appendLoginError(&response, "Ошибка в поле "+err.Key+": "+err.Message, 400)
			log.Println("Error on " + err.Key)
		}
	}

	// Checking for having validation errors
	if response.Errors != nil {
		log.Println("Errors while singing up")
		c.Data["json"] = response
		c.ServeJSON()
		return
	}

	// Getting a user from db
	user := c.getUserByEmail(request.Email)
	if user == nil {
		c.sendError("Bad email or password", 14)
		return
	}

	// Comparing passwords from a client and the database
	var x pk.PasswordHash
	x.Hash = make([]byte, 32)
	x.Salt = make([]byte, 16)
	var err error

	if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
		log.Println("ERROR:", err)
	}
	if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
		log.Println("ERROR:", err)
	}

	if !pk.MatchPassword(request.Password, &x) {
		c.sendError("Bad email or password", 14)
		return
	}

	// Generating a token and sending it to a client
	token := c.generateToken(user.Id)
	response.IdToken = token

	c.Data["json"] = response
	c.ServeJSON()
}

func (c *MainController) SignUp() {
	var request models.SignUpRequest
	var response models.LoginResponse

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Bad request", 400, 400)
		return
	}

	log.Println(request)

	// Validate input fields
	valid := validation.Validation{}
	valid.Email(request.User.Email, "email")
	valid.Required(request.User.Email, "email")
	valid.Required(request.User.Password, "password")
	valid.Required(request.User.Group, "group")
	valid.Required(request.User.FirstName, "first_name")
	valid.Required(request.User.SecondName, "second_name")
	valid.Required(request.User.LastName, "last_name")
	valid.Required(request.User.Gender, "gender")
	valid.MaxSize(request.User.Email, 30, "email")
	valid.MaxSize(request.User.Password, 30, "email")
	valid.MaxSize(request.User.FirstName, 25, "first_name")
	valid.MaxSize(request.User.SecondName, 25, "second_name")
	valid.MaxSize(request.User.LastName, 25, "last_name")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.appendLoginError(&response, "Ошибка в поле "+err.Key+": "+err.Message, 400)
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

	// If it's alright, hashing the password and creating a new user
	h := pk.HashPassword(request.User.Password)
	o := orm.NewOrm()

	request.User.Password = hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt)
	if _, err := o.Insert(&request.User); err != nil {
		c.sendError("Couldn't sign up. The user probably already exists", 14)
		log.Println(err)
	} else {
		// Generating a token and sending it to a client
		token := c.generateToken(request.User.Id)
		response.IdToken = token
	}

	c.Data["json"] = response
	c.ServeJSON()
}

func (c *MainController) generateToken(userId int64) string {
	// Filling the payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "http://localhost:8080",
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + 36000,
	})

	tokenString, _ := token.SignedString(privateKey)

	return tokenString
}

// Checks authorization header
// If it's valid, returns a payload
func (c *MainController) validateToken() (jwt.MapClaims, error) {
	// Getting a token from request header
	tokenString := c.Ctx.Input.Header("Authorization")
	if tokenString == "" {
		log.Println("Token not found")
		return nil, errors.New("Couldn't find Authorization header")
	}
	log.Println(tokenString)

	token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if token == nil {
		return nil, errors.New("Token is null")
	}

	if token.Valid {
		return token.Claims.(jwt.MapClaims), nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			log.Println("That's not even a token")
			return nil, err
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			log.Println("Timing is everything")
			return nil, err
		} else {
			log.Println("Couldn't handle this token:", err)
			return nil, err
		}
	} else {
		log.Println("Couldn't handle this token:", err)
		return nil, err
	}
}

// Checks that json field is string and appends an error into response if it isn't
func (c *MainController) checkStringField(userProperty *string, field interface{}, response *models.LoginResponse, fieldName string) {
	if checkedField, ok := field.(string); ok {
		*userProperty = checkedField
		log.Println(checkedField)
	} else {
		c.appendLoginError(response, "Datatype error in "+fieldName, 400)
	}
}

// The same as previous but for int64
func (c *MainController) checkIntField(userProperty *int64, field interface{}, response *models.LoginResponse, fieldName string) {
	if stringField, ok := field.(string); ok {
		if checkedField, err := strconv.ParseInt(stringField, 10, 64); err != nil {
			c.appendLoginError(response, "Datatype error in "+fieldName, 400)
		} else {
			*userProperty = checkedField
		}
	} else {
		c.appendLoginError(response, "Datatype error in "+fieldName, 400)
	}
}

// Appends an error into the response body
func (c *MainController) appendLoginError(response *models.LoginResponse, message string, code float64) {
	response.Errors = append(response.Errors, models.Error{
		UserMessage: message,
		Code:        code,
	})
}

// The same but for get user response
// Where are my generics Google?!
func (c *MainController) appendGetUserInfoError(response *models.GetUserInfoResponse, message string, code float64) {
	response.Errors = append(response.Errors, models.Error{
		UserMessage: message,
		Code:        code,
	})
}
