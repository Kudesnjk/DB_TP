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
insert into threads (id, created, title, message, user_nickname, forum_slug) values(default, default, $1, $2, $3, $4);

-- name: SelectForumBySlug :one
select f.slug, f.title, f.user_nickname from forums as f join 