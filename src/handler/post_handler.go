package handler

import (
	"alming_backend/src/engine"
	"alming_backend/src/repository"
	"alming_backend/src/service"
	"log"
)

func AddPost(c *engine.Context) {
	post := new(repository.Post)
	c.ParseJSONParam(post)
	service.AddPost(post)
	// service.SavePost(post)

}

func GetPosts(c *engine.Context) {
	posts := service.GetPosts(c)
	c.WriteJSON(posts)
}

func GetPost(c *engine.Context) {
	post := service.GetPost(c)
	c.WriteJSON(post)
}

func UpdatePost(c *engine.Context) {
	post := new(repository.Post)
	c.ParseJSONParam(post)
	service.UpdatePost(post)
}

func DeletePost(c *engine.Context) {
	id, err := c.PathParamToInt("id")
	if err != nil {
		log.Println("Path variable to int fail")
	}
	service.DeletePost(&repository.Post{
		Id: id,
	})
}

func MarkDownUpload(c *engine.Context) {
	service.MarkDownUpload(c)
}
