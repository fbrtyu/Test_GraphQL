package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"ozon-test/graph/model"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type DB struct {
	connect *sql.DB
}

// Подключение к базе данных
func ConnectPG() *DB {
	errr := godotenv.Load()
	if errr != nil {
		log.Fatal("Error loading .env file")
	}
	dbuser, errconn := os.LookupEnv("DBUSER")
	if errconn {
		fmt.Println(dbuser)
	}
	dbpassword, errconn := os.LookupEnv("DBPASSWORD")
	if errconn {
		fmt.Println(dbuser)
	}
	dbname, errconn := os.LookupEnv("DBNAME")
	if errconn {
		fmt.Println(dbuser)
	}
	connStr := "user=" + dbuser + " password=" + dbpassword + " dbname=" + dbname + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return &DB{
		connect: db,
	}
}

// Создание нового пользователя и запись его данных в БД
func (database *DB) CreateUser(input *model.CreateUserInput) *model.User {
	userid := 0
	err := database.connect.QueryRow("insert into users (username) values ($1) returning id", input.Username).Scan(&userid)
	if err != nil {
		panic(err)
	}
	return &model.User{
		ID: &userid,
	}
}

// Создание нового поста и запись в БД
func (database *DB) CreatePost(userid int, input *model.CreatePostInput) *model.Post {
	postid := 0
	err := database.connect.QueryRow("insert into post (title, text, commenting, iduser) values ($1, $2, $3, $4) returning id", input.Title, input.Text, input.Commenting, userid).Scan(&postid)
	if err != nil {
		panic(err)
	}
	return &model.Post{
		ID: &postid,
	}
}

// Создание комментария к посту и запись в БД
func (database *DB) CreateComment(userid int, postid int, input *model.CreateCommentInput) *model.Comment {
	commentid := 0
	err := database.connect.QueryRow("insert into comment (idpost, text, iduser) values ($1, $2, $3) returning id", postid, input.Text, userid).Scan(&commentid)
	if err != nil {
		panic(err)
	}
	return &model.Comment{
		ID: &commentid,
	}
}

// Создание комментария, который будет являться ответом на уже существующий комментарий. Запись данных в БД
func (database *DB) CreateAnswer(postid int, commentid int, answerid int, input *model.CreateCommentInput) *model.Comment {
	err := database.connect.QueryRow("insert into commentanswer (idcomment, idanswer, idpost) values ($1, $2, $3)", commentid, answerid, postid)
	if err != nil {
		panic(err)
	}
	err = database.connect.QueryRow("update comment set answered = 'true' where id = $1", commentid)
	if err != nil {
		panic(err)
	}
	return &model.Comment{
		ID: &commentid,
	}
}

// Получение данных из БД, список постов
func (database *DB) Post() []*model.Post {
	rows, err := database.connect.Query("select * from post")
	if err != nil {
		panic(err)
	}
	var posts = []*model.Post{}
	for rows.Next() {
		var post = model.Post{}
		var user = model.User{}
		post.User = &user
		err := rows.Scan(&post.ID, &post.Title, &post.Text, &post.Commenting, &post.User.ID)
		if err != nil {
			fmt.Println(err)
		}
		posts = append(posts, &post)
	}
	return posts
}

// Получение данных из БД для вывода определенного поста и всех комментариев
func (database *DB) PostAndComments(id int) *model.Post {
	var postwithcomment = model.Post{}
	var userinfo = model.User{}
	postwithcomment.User = &userinfo

	//Получаем из БД информацию о посте
	postinf, err := database.connect.Query("select id, title, text, commenting, iduser from post where id = $1", id)
	if err != nil {
		panic(err)
	}
	for postinf.Next() {
		err := postinf.Scan(&postwithcomment.ID, &postwithcomment.Title, &postwithcomment.Text, &postwithcomment.Commenting, &postwithcomment.User.ID)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Получаем информацию о комментариях, на которые нет ответа
	rowsOfUncomment, err := database.connect.Query("select id, idpost, text, iduser from comment where idpost = $1 and answered = 'false' order by id", id)
	if err != nil {
		panic(err)
	}
	for rowsOfUncomment.Next() {
		var comment = model.Comment{}
		comment.User = &userinfo
		err := rowsOfUncomment.Scan(&comment.ID, &comment.Postid, &comment.Text, &comment.User.ID)
		if err != nil {
			fmt.Println(err)
		}

		// Добавляем эти комментарии первыми в массив комментариев поста
		postwithcomment.Comments = append(postwithcomment.Comments, &comment)
	}

	// Далее была попытка за один проход по полученной таблицы из БД записать комментарии и ответы на них в структуры
	// Вывести в консоль друг на против друга комментарий и ответ на него я могу, но есть проблемы как эти данные правильно сформировать для корректного ответа клиенту
	// ПРИ ХРАНЕНИИ ДАННЫХ В ПАМЯТИ ПРОГРАММЫ, РЕАЛИЗАЦИЯ И РАБОТА ВЫВОДА ВСЕХ КОММЕНТАРИЕВ ПОСТА ГОРАЗДО ЛУЧШЕ!
	var answerarr = []*model.Comment{}
	var answerarrlast = []*model.Comment{}
	var currentid = 1
	var lastid = 0
	var flag = true

	// Получение двнных, которые содержат информаию о комментарии и всех ответах, которые есть на него
	rowsOfAnsComment, err := database.connect.Query(`select A.id, A.text, A.iduser as iduser1, A.idanswer, B.text, B.iduser as iduser2 from 
													(SELECT comment.id, comment.idpost, text, comment.iduser, commentanswer.idanswer FROM comment, commentanswer WHERE comment.id = commentanswer.idcomment AND comment.idpost = $1) A 
													join (SELECT id, text, iduser FROM comment where idpost = $2) B on A.idanswer = B.id order by id`, id, id)
	if err != nil {
		panic(err)
	}

	// Переменные для хранения текущего комментария и предыдущего
	var comment = model.Comment{}
	comment.User = &userinfo
	var commentlast = model.Comment{}
	commentlast.User = &userinfo

	for rowsOfAnsComment.Next() {
		// Тут необходимо было пересоздавать структуру для хранения нового комментария. Перед этим запомнив его
		var answer = model.Comment{}
		answer.User = &userinfo
		comment.Answer = answerarr
		commentlast.Answer = answerarrlast

		err := rowsOfAnsComment.Scan(&comment.ID, &comment.Text, &comment.User.ID, &answer.ID, &answer.Text, &answer.User.ID)
		if err != nil {
			fmt.Println(err)
		}

		currentid = *comment.ID

		fmt.Println(currentid)
		fmt.Println(lastid)
		fmt.Println(flag)

		// Основная логика такова:
		// Берутся данные первой строки. Далее идёт проверка совпадения текущего id комментария и прошлого.
		// Первый раз оно специально совпадает, чтобы не пропустить запись (1 == 1)
		// Когда уже пошли на вторую итерацию, то из-за flag первое условие уже никогда не будет выполняться

		// Теперь на каждой итерации всё так же проверяются id комментариев
		// Когда они совпадают, то массив ответов на текущий комментарий копится (так как в таблице, полученной ранее из БД, записи ответов на один и тот же комментарий идут друг за другом)
		// Когда они не совпадают (то есть начались ответы на уже другой комментарий), происходит запись данного массив ответов в массив ответов комментария
		// И данный комментарий с ответами на него записывается в массив комментариев уже самого поста. Данные такого поста возвращаются в ответ
		if currentid == *comment.ID && flag {
			answerarr = append(answerarr, &answer)
			lastid = *comment.ID
			flag = false
			fmt.Println(comment, " - ", answer)
			commentlast = comment
			answerarrlast = answerarr
		} else if currentid == lastid {
			fmt.Println(comment, " - ", answer)
			answerarr = append(answerarr, &answer)
			lastid = *comment.ID
		} else if currentid != lastid {
			commentlast.Answer = append(commentlast.Answer, answerarrlast...)
			postwithcomment.Comments = append(postwithcomment.Comments, &commentlast)
			fmt.Println(comment, " - ", answer)
			comment.Answer = append(comment.Answer, answerarr...)
			postwithcomment.Comments = append(postwithcomment.Comments, &comment)
			comment.Answer = []*model.Comment{}
			comment.User = &userinfo
			commentlast.Answer = []*model.Comment{}
			commentlast.User = &userinfo
			lastid = *comment.ID
		}
	}
	return &postwithcomment
}
