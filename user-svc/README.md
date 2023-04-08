# Authentication Service
This is a microservice that handle all incoming http request and communicate with other microservice on behalf.

# Basic Authentication
Path: /basic-auth

Method: Post

Params: 
- email
- password

# Create User
Path: /user

Method: PUT

Params:
- email: string
- password: string