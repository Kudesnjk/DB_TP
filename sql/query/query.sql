-- name: InsertUser :exec
insert into users (nickname, email, fullname, about) values($1, $2, $3, $4);


-- name: SelectUserByNickname :one
select fullname, nickname, email, about from users where nickname = $1;

-- name: UpdateUser :exec
update users
   set about = $1 , email = $2, fullname = $3 
 where nickname = $4;

-- name: InsertForum :exec
insert into forums (slug, title, user_nickname) values($1, $2, $3); 

-- name: InsertThread :exec
insert into threads (title, message, user_nickname) values($1, $2, $3);