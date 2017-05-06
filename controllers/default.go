package controllers

import (
	"github.com/astaxie/beego"
	"log"
	"os"
	"activist_api/models"
	"strconv"
	//"gopkg.in/gographics/imagick.v2/imagick"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.sendErrorWithStatus("Page not found", 404, 404)
	c.ServeJSON()
}

func (c *MainController) IndexPage() {
	var response models.IndexPageResponse
	response.SoonerEvents = models.GetSoonerEvents(3)
	tags := models.GetTopFiveTags()
	response.EventsByTags = models.GetTopFiveEventsByTags(tags)
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *MainController) GenerateTemplateToken() {

	tplId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		log.Fatal(err)
		return
	}

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		token := generateFileToken(user.Id, tplId, "tpl")
		tpl := models.GetTemplateById(tplId)
		if tpl == nil {
			c.sendErrorWithStatus("Template not found", 404, 404)
			return
		}
		response := models.GenerateTemplateTokenResponse{ Ok: true, Token: token, Template: tpl }
		c.Data["json"] = &response
		c.ServeJSON()
	}
}

func (c *MainController) GenerateFormToken() {
	defer c.ServeJSON()

	formId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		log.Fatal(err)
		return
	}

	if payload, err := validateToken(c.Ctx.Input.Header("Authorization")); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := models.GetUserById(int64(payload["sub"].(float64)))
		if models.IsAllowedToDownloadForm(user.Id, formId) {
			token := generateFileToken(user.Id, formId, "form")
			form := models.GetFormUserById(formId)
			if form == nil {
				c.sendErrorWithStatus("Form not found", 404, 404)
				return
			}
			response := models.GenerateFormTokenResponse{ Ok: true, Token: token, Form: form }
			c.Data["json"] = &response
		} else {
			c.sendErrorWithStatus("Access denied", 403, 403)
		}
	}
}

// Checking a token to give a link to a file.
// I get a token in the get request, and that's not the best way,
// because a user can send a link to another person and he would give him his token.
// But I don't know how to do that another way, so I gotta fix it later.
func (c *MainController) XAccelTemplate() {
	defer c.ServeJSON()

	// Get organizer id and template name from the request
	orgId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		log.Fatal(err)
		return
	}
	path := c.Ctx.Input.Param(":name")


	if payload, err := validateFileToken(c.Input().Get("token")); err != nil {
		log.Println(err)
		c.Ctx.Output.Header("X-Accel-Redirect", "/")
		return
	} else {
		tpl := models.GetTemplateByOrgAndName(orgId, path)
		if tpl == nil {
			log.Println("Template not found")
			c.Ctx.Output.Header("X-Accel-Redirect", "/unauthorized")
			return
		}
		// Check if the token was generated for that file
		if payload["typ"] != "tpl" || int64(payload["fid"].(float64)) != tpl.Id {
			log.Println("Wrong token")
			c.Ctx.Output.Header("X-Accel-Redirect", "/forbidden")
			return
		}
	}
	c.Ctx.Output.Header("X-Accel-Redirect", "/api/storage/docs/tpl/" + strconv.Itoa(int(orgId)) + "/" + path)
}

func (c *MainController) XAccelForm() {
	defer c.ServeJSON()
	var userId int64
	// Get form id from the :id param
	prtId, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 0, 64)
	if err != nil {
		log.Fatal(err)
		return
	}
	path := c.Ctx.Input.Param(":name")

	form := models.GetFormByPrtAndName(prtId, path)
	if form == nil {
		log.Println("Form not found")
		c.Ctx.Output.Header("X-Accel-Redirect", "/unauthorized")
		return
	}

	// Token validation
	if payload, err := validateFileToken(c.Input().Get("token")); err != nil {
		// If the token is invalid
		log.Println(err)
		c.Ctx.Output.Header("X-Accel-Redirect", "/unauthorized")
		return
	} else {
		// Check "typ" field in the token. It should be for volunteer forms
		if payload["typ"] != "form" {
			c.Ctx.Output.Header("X-Accel-Redirect", "/forbidden")
			return
		}
		user := models.GetUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	// Only form and template owners are allowed to download the file
	if !models.IsAllowedToDownloadForm(userId, form.Id) {
		c.Ctx.Output.Header("X-Accel-Redirect", "/forbidden")
		return
	}

	c.Ctx.Output.Header("X-Accel-Redirect", "/api/storage/docs/form/" + strconv.Itoa(int(form.ParticipantId)) + "/" + form.Path)
}

func (c *MainController) UploadFile() {/*
	log.Println("Uploading...")
	file, header, _ := c.GetFile("file") // where <<this>> is the controller and <<file>> the id of your form field
	if file != nil {
		// get the filename
		fileName := header.Filename
		// Get a file extension
		ext := fileName[strings.LastIndex(fileName, "."):]
		// Make a random md5 name for an image
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
		path += "/" + newName + ext
		err := c.SaveToFile("file", path)
		log.Println(err)

		var sFile *os.File
		if sFile, err = os.Open(path); err != nil {
			log.Println(err)
			c.sendError("Couldn't open a file", 1)
			return
		}
		if ok := transformImage(sFile, path); !ok {
			c.sendError("Couldn't transform an image", 1)
		}
		c.sendSuccess()
	} else {
		c.sendError("Meh", 1)
	}*/
}


func transformCover(file *os.File, path string) (bool) {
/*
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.ReadImageFile(file); err != nil {
		log.Println(err)
		return false
	}

	// Use fixed aspect ratio
	const ratio = 1.9

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	if float64(width)/float64(height) < ratio {
		newHeight := uint(float64(width) / ratio)
		deltay := (height - newHeight) / 2
		if err := mw.CropImage(width, newHeight, 0, int(deltay)); err != nil {
			log.Println(err)
			return false
		}
	} else if float64(width)/float64(height) > ratio {
		newWidth := uint(float64(height) * ratio)
		deltax := (width - newWidth) / 2
		if err := mw.CropImage(newWidth, height, int(deltax), 0); err != nil {
			log.Println(err)
			return false
		}
	}

	if err := mw.ResizeImage(1140, 600, imagick.FILTER_LANCZOS, 1); err != nil {
		log.Println(err)
		return false
	}
	if err := mw.SetImageCompression(imagick.COMPRESSION_JPEG); err != nil {
		log.Println(err)
		return false
	}
	if err := mw.StripImage(); err != nil {
		log.Println(err)
		return false
	}
	if err := mw.SetImageCompressionQuality(90); err != nil {
		log.Println(err)
		return false
	}
	if err := mw.WriteImage(path); err != nil {
		log.Println(err)
		return false
	}*/
	return true
}

func transformAvatar(file *os.File, path string) (bool) {
	/*
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.ReadImageFile(file); err != nil {
		log.Println(err)
		return false
	}

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	if width < height {
		deltay := (height - width) / 2
		if err := mw.CropImage(width, width, 0, int(deltay)); err != nil {
			log.Println(err)
			return false
		}
	} else if width > height {
		deltax := (width - height) / 2
		if err := mw.CropImage(height, height, int(deltax), 0); err != nil {
			log.Println(err)
			return false
		}
	}

	if err := mw.ResizeImage(250, 250, imagick.FILTER_LANCZOS, 1); err != nil {
		log.Println(err)
		return false
	}
	if err := mw.SetImageCompression(imagick.COMPRESSION_JPEG); err != nil {
		log.Println(err)
		return false
	}
	if err := mw.StripImage(); err != nil {
		log.Println(err)
		return false
	}
	if err := mw.SetImageCompressionQuality(90); err != nil {
		log.Println(err)
		return false
	}
	if err := mw.WriteImage(path); err != nil {
		log.Println(err)
		return false
	}*/
	return true
}
