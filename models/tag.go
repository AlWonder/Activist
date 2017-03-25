package models

import (
	"github.com/astaxie/beego/orm"
	"log"
  "strings"
)

type Tag struct {
	Id   int64  `orm:"column(id);auto" json:"id"`
	Name string `orm:"column(name);size(30);unique" json:"name,omitempty"`
}

func GetTags(tag string) *[]string {
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

func GetTagStatus(userId, tagId int64) *TagStatus {
	var tagStatus TagStatus
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

func AddTagStatus(name string, userId int64, status bool) int64 {
	o := orm.NewOrm()
	tagStatus := TagStatus{UserId: userId, Status: status}
	if tag := GetTag(name); tag != nil {
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

func DeleteTagStatus(tagId, userId int64) error {
	o := orm.NewOrm()
	if _, err := o.Raw(`DELETE
			FROM tags_status
			WHERE tag_id = ? AND user_id = ?`, tagId, userId).Exec(); err != nil {
		return err
	}
	return nil
}

func GetTag(name string) *Tag {
	o := orm.NewOrm()
	log.Println(name)

	tag := Tag{Name: name}

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

func AddTags(tags []string) []int64 {
	var tagIds []int64
	o := orm.NewOrm()
	for _, tag := range tags {
		newTag := Tag{
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

func DeleteEventTags(eventId int64, tags []string) error {
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

func AddEventTags(eventId int64, tagIds []int64) bool {
	var ok bool
	o := orm.NewOrm()
	for _, tagId := range tagIds {
		event := Event{Id: eventId}
		tagEvent := EventTag{
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

func GetTagsByEventId(id int64) *[]string {
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

func GetTopFiveTags() *[]Tag {
	var tags []Tag
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
