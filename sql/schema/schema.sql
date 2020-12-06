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
    foreign key (user_nickname) references users(nickname) on delete cascade
);

create table threads (
    id serial primary key,
    title text not null,
    message text not null,
    created timestamp default now(),
    user_nickname text not null,
    forum_id text not null,
    foreign key (user_nickname) references users(nickname) on delete cascade ,
    foreign key (forum_id) references forums(slug) on delete cascade
);

create table posts (
    id serial primary key ,
    message text not null ,
    is_edited boolean default false,
    created timestamp default now(),
    parent_id int default 0,
    user_nickname text not null ,
    thread_id int not null,
    foreign key (user_nickname) references users(nickname) on delete cascade ,
    foreign key (thread_id) references threads(id) on delete cascade ,
    foreign key (parent_id) references posts(id) on delete cascade
);
