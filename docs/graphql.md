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
