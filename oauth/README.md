# OAuth Flow
Go to github, and create a Client ID and Genrate a new Client Secret.
This tells Github OAuth server that, there is an app with id `ClientID`, that might ask you to authorize things.

## The Flow
* User clicks `/login`
* Server redirects with `github.com/oauth/authorize?client_id=a&scope=user&redirect_uri=localhost/callback`
* The OAuth server sends a `code` to server by calling `localhost/callback?code=1234`
* Server requests oauth server for access token by sending `secret` and `code`
* Server uses the access token to access the account lol
