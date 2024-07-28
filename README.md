# Ozon_test

<h2>Запросы GraphQL:</h2>
<pre>
mutation addUser {
  createUser(input: {username: "newuser"}) {
    id
  }
}
</pre>
<pre>
mutation addPost {
  createPost(
    userid: 1
    input: {title: "nameofpost", text: "sometext", commenting: true}
  ) {
    id
  }
}
</pre>
<pre>
mutation addComment {
  createComment(userid: 1, postid: 1, input: {text: "somecomment"}) {
    id
  }
}
</pre>
<pre>
mutation addAnswer {
  createAnswer(userid: 1, postid: 1, commentid: 1, input: {text: "answer"}) {
    id
  }
}
</pre>
<pre>
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
</pre>
<pre>
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
</pre>
<pre>
subscription Subscription {
  comment {
    text
  }
}
</pre>
