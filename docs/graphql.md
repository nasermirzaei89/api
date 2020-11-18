# GraphQL

## Log In

```graphql
mutation ($username: String!, $password: String!) {
    logIn(request: {username: $username, password: $password}) {
        accessToken
        user {
            uuid
            username
        }
    }
}
```

## Get Me (Current User)

```graphql
query  {
    me {
        uuid
        username
    }
}
```

## Create a Post

```graphql
mutation ($title: String!, $slug: String!, $contentMarkdown: String!) {
    createPost(request: {title: $title, slug: $slug, contentMarkdown: $contentMarkdown}) {
        uuid
        title
        slug
        contentMarkdown
        contentHTML
    }
}
```
