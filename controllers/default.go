package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/gographics/imagick.v2/imagick"
	"log"
	"os"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

/*func (c *MainController) UploadFile() {
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
	}
}*/

func transformCover(file *os.File, path string) (bool) {
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
	}
	return true
}

func transformAvatar(file *os.File, path string) (bool) {
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
	}
	return true
}
