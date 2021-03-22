package service

import (
	"alming_backend/src/engine"
	"alming_backend/src/repository"
	"io"
)

func GetPosts(c *engine.Context) (posts []repository.Post) {
	repository.GetPosts(&posts)
	return posts
}

func GetPost(c *engine.Context) *repository.Post {
	id, err := c.PathParamToInt("id")
	if err == nil {
		return repository.GetPost(id)
	}
	return nil
}

func MarkDownUpload(c *engine.Context) {
	title := c.PostParam("title")
	md, _, _ := c.File("post")
	buf := make([]byte, 128)
	var content []byte
	for {
		readLen, err := md.Read(buf)
		if readLen == 0 && err == io.EOF {
			break
		}
		content = append(content, buf[0:readLen]...)
	}

	SavePost(repository.CreateMdPostDefault(title, content))
}

func SavePost(post repository.Post) {
	repository.InsertPost(&post)
}
