package localdatabase

import (
	"ozon-test/graph/model"
	"ozon-test/subscription"
)

// Переменные для хранения данных в мапяти программы
var users []*model.User
var posts []*model.Post
var comments []*model.Comment

// Счётчики для id, чтобы они всегда были уникальными внутри одного массива структур
var usersid = 0
var postsid = 0
var commentsid = 0

// Используется при создании нового комментария
// Если передавать значение ommentsid в ID при создании новго комментария, то значения id будут некорректными, так как значение ommentsid успеет измениться
var arrayofcommentsid []int

// В данном случае всё почти аналогично работе в database.go, но проще так как, например, можно сразу обратиться к нужному элементу в массиве, а не идти по таблице

func CreateUser(input *model.CreateUserInput) *model.User {
	// Изменение счётчика id
	usersid = usersid + 1

	// Создание нового пользователя
	newuser := model.User{
		ID:       &usersid,
		Username: input.Username,
	}

	// Добавление в массив пользователей
	users = append(users, &newuser)

	// Возвращаем как и ранее, id пользователя
	return &model.User{
		ID: &usersid,
	}
}

func CreatePost(userid int, input *model.CreatePostInput) *model.Post {
	postsid = postsid + 1

	// Тут необходимо использовать (User: users[userid-1]) так как нумерация в массиве начинается с 0, а значения id я начинаю с 1. То есть если вычесть 1 из id записи, можно сразу узнать её индекс в массиве
	newpost := model.Post{
		ID:         &postsid,
		User:       users[userid-1],
		Title:      input.Title,
		Text:       input.Text,
		Commenting: input.Commenting,
	}
	posts = append(posts, &newpost)
	return &model.Post{
		ID: &postsid,
	}
}

func CreateComment(userid int, postid int, input *model.CreateCommentInput) *model.Comment {
	commentsid = commentsid + 1
	arrayofcommentsid = append(arrayofcommentsid, commentsid)
	newcomment := model.Comment{
		ID:     &arrayofcommentsid[commentsid-1],
		User:   users[userid-1],
		Postid: postid,
		Text:   input.Text,
	}
	comments = append(comments, &newcomment)
	posts[postid-1].Comments = append(posts[postid-1].Comments, comments[commentsid-1])
	newcommenta := subscription.NewComment{}
	newcommenta.Text = input.Text
	subscription.Publicch <- newcommenta
	return &model.Comment{
		ID: &commentsid,
	}
}

func CreateAnswer(userid int, postid int, commentid int, input *model.CreateCommentInput) *model.Comment {
	commentsid = commentsid + 1
	arrayofcommentsid = append(arrayofcommentsid, commentsid)
	newcomment := model.Comment{
		ID:     &arrayofcommentsid[commentsid-1],
		User:   users[userid-1],
		Postid: postid,
		Text:   input.Text,
	}

	// Добавляем новый ответ на комментарий в массив комментариев
	// Далее добавлем этот ответ в массив Answer комментария, на который и давался этот ответ
	// Каждый Post ссылается на Commet, который содержит Answer, который в свою очередь содержит Comment и т. д.
	// Таким образом работает вложенность ответов на комментарии и ответов на ответы и т. д.
	// Но я не совсем понял как в браузере при GraphQL запросе вывести данную вложенность до бесконечности
	comments = append(comments, &newcomment)
	comments[commentid-1].Answer = append(comments[commentid-1].Answer, &newcomment)
	return &model.Comment{
		ID: &commentsid,
	}
}

func Posts() []*model.Post {
	var arrayposts = []*model.Post{}
	for i := 0; i < len(posts); i++ {
		arrayposts = append(arrayposts, posts[i])
	}
	return arrayposts
}

func PostAndComments(id int) *model.Post {
	return posts[id-1]
}
