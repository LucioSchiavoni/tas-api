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


Ver notificaciones que le llegan al usuario:
GET /notificationByUser/1

Esto devuelve:
ID: (id de notificacion)
CreatedAt
UserID: (id del que recibe esta notificacion)
Type: si es de tipo comments o like
PostID: id de ese post creado
Check: por defecto viene en falso pero al verlo debe cambiar a true 
CreatorID: id del usuario que hizo la notificacion (persona que hizo el like o creo el comentario)


CHAT:

Crear mensaje:

POST: 
http://localhost:8080/message?sender={1}&recipient={4}

{
  "Sender": "{1}",
  "Recipient": "{2}",
  "Content": "Hola, ¿cómo estás?"
}

Falta el get message 