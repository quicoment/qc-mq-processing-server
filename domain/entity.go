package domain

import (
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	Password  string    `json:"password"`
	Likes     int       `json:"likes"`
}
//
//func insert(c Comment) error {
//	data, _ := json.Marshal(c)
//	err := common.POST(data)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	return nil
//}
//
//func update(c Comment) error {
//	data, _ := json.Marshal(c)
//	err := common.UPDATE(c.ID, data)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	return nil
//}
//
//func get(key int64) Comment {
//	result, err := common.Get(key)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//
//	var comment Comment
//	json.Unmarshal(result, &comment)
//	return comment
//}
//
//func delete(c Comment) error {
//	err := common.Delete(c.ID)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	return nil
//}
