# Ozon_test

<h2>Запросы GraphQL:</h2>
<pre>
 
mutation addUser {
  createUser(input: {username: "newuser"}) {
    id
  }
}

mutation addPost {
  createPost(
    userid: 1
    input: {title: "nameofpost", text: "sometext", commenting: true}
  ) {
    id
  }
}

mutation addComment {
  createComment(userid: 1, postid: 1, input: {text: "somecomment"}) {
    id
  }
}

mutation addAnswer {
  createAnswer(userid: 1, postid: 1, commentid: 1, input: {text: "answer"}) {
    id
  }
}

query GetListOfPosts {
  posts {
    id
    user {
      id
    }
    title
    text
    commenting
  }
}

query GetPostAndComments {
  post(id: 1) {
    id
    user {
      id
    }
    title
    text
    commenting
    comments {
      id
      user {
        id
      }
      text
      answer {
        id
        user {
          id
        }
        text
        answer {
          id
          user {
            id
          }
          text
        }
      }
    }
  }
}

subscription Subscription {
  comment {
    text
  }
}

</pre>

<h2>Структура БД</h2>
https://raw.githubusercontent.com/fbrtyu/Ozon_test/main/Untitled.jpg

<h2>Демонстрация работы</h2>
https://youtu.be/ywaus5FWiIs?si=GXwykTv0UQFzzv6_
