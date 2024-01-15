# tas-api

Endpoints:

User:
Create user:  POST /user (formData)

Login: POST /login (json)
Authentication token: GET /auth Authorization Bearer token

Post: 
Create post: POST /post

Post by user id: GET /post/user_id


Like:

Crear like: POST /like  
ejemplo 
raw:
{
    "UserID":2,  //usuario que da el like
    "PostID":1   // id del post
}