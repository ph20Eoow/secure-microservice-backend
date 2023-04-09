# Authentication Service
This is a microservice that handle all incoming http request and communicate with other microservice on behalf.

| method | path           | requireAuth | description          |   |
|--------|----------------|-------------|----------------------|---|
| get    | /user/{id}     | true        | update user          |   |
| put    | /user/         | false       | create user          |   |
| post   | /user/password | treu        | update user password |   |