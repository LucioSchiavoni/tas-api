# tas-api

Endpoints:

User:
Create user:  POST /user (formData)

Login: POST /login (json)
Authentication token: GET /auth Authorization Bearer token

Post: 
Create post: POST /post

Post by user id: GET /post/user_id

Todos los Post: GET /AllPost  


Like:

Crear like: POST /like  
ejemplo 
raw:
{
    "UserID":2,  //usuario que creo el post
    "PostID":1   // id del post
    "CreatorID":  //usuario que da el like
}

Comentarios: 

Crear comentario: POST /comments
raw:
{   
    "UserID":
    "PostID":
    "content":"hola"
    "CreatorID" : 
}

Commentarios por id del post:
GET  /{post_id}/comments

