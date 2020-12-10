create table users (
    nickname text primary key,
    email text unique not null,
    fullname text not null ,
    about text
);

create table forums (
    slug text primary key ,
    title text not null,
    user_nickname text not null,
    threads_num int not null default 0,
    posts_num int not null default 0,
    foreign key (user_nickname) references users(nickname) on delete cascade
);

create table threads (
    id serial primary key,
    title text not null,
    message text not null,
    created timestamp not null default now(),
    user_nickname text not null,
    forum_slug text not null,
    votes int not null default 0,
    foreign key (user_nickname) references users(nickname) on delete cascade ,
    foreign key (forum_slug) references forums(slug) on delete cascade
);

create table posts (
    id serial primary key ,
    message text not null ,
    is_edited boolean not null default false,
    created timestamp not null default now(),
    parent_id int not null default 0,
    user_nickname text not null ,
    thread_id int not null,
    foreign key (user_nickname) references users(nickname) on delete cascade ,
    foreign key (thread_id) references threads(id) on delete cascade ,
    foreign key (parent_id) references posts(id) on delete cascade
);


GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO db_forum_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO db_forum_user;