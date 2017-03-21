package controllers

import (
	"activist_api/models"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"log"
	"strings"
)

func (c *MainController) QueryTags() {
	tag := c.Input().Get("query")
	tags := c.getTags(tag)
	c.Data["json"] = &tags
	c.ServeJSON()
}

func (c *MainController) GetTagStatus() {
	var tagName string
	var userId, tagId int64
	tagName = c.Ctx.Input.Param(":tag")

	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := c.getUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	if tag := c.getTag(tagName); tag == nil {
		c.sendError("Tag not found", 14)
		return
	} else {
		tagId = tag.Id
	}

	var response models.GetTagStatusResponse

	if tagStatus := c.getTagStatus(userId, tagId); tagStatus == nil {
		response.HasStatus = false
	} else {
		response.HasStatus = true
		response.Status = tagStatus.Status
	}

	response.Ok = true
	c.Data["json"] = &response
	c.ServeJSON()
}

func (c *MainController) DeleteTagStatus() {
	var tagName string
	var userId, tagId int64
	tagName = c.Ctx.Input.Param(":tag")

	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := c.getUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	if tag := c.getTag(tagName); tag == nil {
		c.sendError("Tag not found", 14)
		return
	} else {
		tagId = tag.Id
	}

	if err := c.deleteTagStatus(tagId, userId); err != nil {
		c.sendError("Couldn't delete tag status", 14)
		return
	}

	c.sendSuccess()
}

func (c *MainController) AddTagStatus() {
	var tag string
	var userId int64
	var status bool
	tag = c.Ctx.Input.Param(":tag")

	if payload, err := c.validateToken(); err != nil {
		log.Println(err)
		c.sendErrorWithStatus("Invalid token. Access denied", 401, 401)
		return
	} else {
		user := c.getUserById(int64(payload["sub"].(float64)))
		userId = user.Id
	}

	var request models.AddFavHideTagRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err == nil {
		status = request.Status
	}

	if statusId := c.addTagStatus(tag, userId, status); statusId == 0 {
		c.sendError("Couldn't add tag status", 14)
		return
	}
	c.sendSuccess()
}

func (c *MainController) getTags(tag string) *[]string {
	var tags []string

	o := orm.NewOrm()
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}

	like := "%" + tag + "%"

	_, err := o.Raw(`SELECT name
					 FROM tags
					 WHERE name LIKE ?`, like).QueryRows(&tags)
	if err != nil {
		log.Println("findTags: ", err)
		return nil
	}
	return &tags
}

func (c *MainController) getTagStatus(userId, tagId int64) *models.TagStatus {
	var tagStatus models.TagStatus
	o := orm.NewOrm()

	err := o.Raw(`SELECT *
		FROM tags_status
		WHERE user_id = ? AND tag_id = ?`, userId, tagId).QueryRow(&tagStatus)
	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return nil
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return nil
	}
	return &tagStatus
}

func (c *MainController) addTagStatus(name string, userId int64, status bool) int64 {
	o := orm.NewOrm()
	tagStatus := models.TagStatus{UserId: userId, Status: status}
	if tag := c.getTag(name); tag != nil {
		tagStatus.TagId = tag.Id
	} else {
		return 0
	}

	if created, id, err := o.ReadOrCreate(&tagStatus, "UserId", "TagId", "Status"); err == nil {
		if created {
			log.Println("New Insert an object. Id:", id)
		} else {
			log.Println("Get an object. Id:", id)
		}
		return id
	} else {
		log.Println(err)
		return 0
	}
}

func (c *MainController) deleteTagStatus(tagId, userId int64) error {
	o := orm.NewOrm()
	if _, err := o.Raw(`DELETE
			FROM tags_status
			WHERE tag_id = ? AND user_id = ?`, tagId, userId).Exec(); err != nil {
		return err
	}
	return nil
}

func (c *MainController) getTag(name string) *models.Tag {
	o := orm.NewOrm()
	log.Println(name)

	tag := models.Tag{Name: name}

	err := o.Raw(`SELECT *
		FROM tags
		WHERE name = ?`, name).QueryRow(&tag)

	if err == orm.ErrNoRows {
		log.Println("No result found.")
		return nil
	} else if err == orm.ErrMissPK {
		log.Println("No primary key found.")
		return nil
	}
	return &tag
}

func (c *MainController) addTags(tags []string) []int64 {
	var tagIds []int64
	o := orm.NewOrm()
	for _, tag := range tags {
		newTag := models.Tag{
			Name: tag,
		}

		if created, id, err := o.ReadOrCreate(&newTag, "Name"); err == nil {
			if created {
				log.Println("New Insert an object. Id:", id)
			} else {
				log.Println("Get an object. Id:", id)
			}
			tagIds = append(tagIds, id)
		} else {
			log.Println(err)
		}
	}

	return tagIds
}

func (c *MainController) deleteEventTags(eventId int64, tags []string) error {
	o := orm.NewOrm()
	for _, tag := range tags {
		log.Println("Deleting " + tag)
		if _, err := o.Raw(`DELETE e.*
			FROM events_tags e
			INNER JOIN tags t ON t.id = e.tag_id
			WHERE t.name = ? AND e.event_id = ?`, tag, eventId).Exec(); err != nil {
			return err
		}
	}
	return nil
}

func (c *MainController) addEventTags(eventId int64, tagIds []int64) bool {
	var ok bool
	o := orm.NewOrm()
	for _, tagId := range tagIds {
		event := models.Event{Id: eventId}
		tagEvent := models.EventTag{
			Event: &event,
			TagId: tagId,
		}

		if _, err := o.Insert(&tagEvent); err != nil {
			log.Println(err)
			ok = false
		}
	}
	ok = true
	return ok
}

func (c *MainController) getTagsByEventId(id int64) *[]string {
	var tags []string
	o := orm.NewOrm()
	if _, err := o.Raw(`SELECT name
		FROM tags
		INNER JOIN events_tags ON tags.id = events_tags.tag_id
		WHERE event_id = ?`, id).QueryRows(&tags); err != nil {
		log.Println(err)
		return nil
	}
	return &tags
}

func (c *MainController) getTopFiveTags() *[]models.Tag {
	var tags []models.Tag
	o := orm.NewOrm()
	if _, err := o.Raw(`SELECT t.*
FROM tags t INNER JOIN (events_tags et INNER JOIN events e ON e.id = et.event_id) ON et.tag_id = t.id
GROUP BY t.id
ORDER BY count(*) DESC
LIMIT 5`).QueryRows(&tags); err != nil {
		log.Println(err)
		return nil
	}
	return &tags
}
