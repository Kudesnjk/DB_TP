CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
    id serial,
    nickname citext COLLATE "C" PRIMARY KEY,
    email text UNIQUE NOT NULL,
    fullname text NOT NULL,
    about text
);

CREATE TABLE forums (
    slug text PRIMARY KEY,
    title text NOT NULL,
    user_nickname text NOT NULL,
    threads_num int NOT NULL DEFAULT 0,
    posts_num int NOT NULL DEFAULT 0,
    FOREIGN KEY (user_nickname) REFERENCES users(nickname) ON DELETE CASCADE
);

CREATE TABLE threads (
    id serial PRIMARY KEY,
    slug text NOT NULL,
    title text NOT NULL,
    message text NOT NULL,
    created timestamptz NOT NULL DEFAULT NOW(),
    user_nickname text NOT NULL,
    forum_slug text NOT NULL,
    votes int NOT NULL DEFAULT 0,
    FOREIGN KEY (user_nickname) REFERENCES users(nickname) ON DELETE CASCADE,
    FOREIGN KEY (forum_slug) REFERENCES forums(slug) ON DELETE CASCADE
);

CREATE TABLE votes (
    user_nickname text NOT NULL REFERENCES users (nickname) ON DELETE CASCADE,
    thread_id int NOT NULL REFERENCES threads (id) ON DELETE CASCADE,
    voice int NOT NULL,
    PRIMARY KEY(user_nickname, thread_id)
);

CREATE TABLE posts (
    id serial PRIMARY KEY,
    message text NOT NULL,
    is_edited boolean NOT NULL DEFAULT false,
    created timestamptz NOT NULL DEFAULT NOW(),
    parent_id int NOT NULL DEFAULT 0,
    path integer [] NOT NULL DEFAULT '{}',
    user_nickname text NOT NULL,
    forum_slug text NOT NULL,
    thread_id int NOT NULL,
    FOREIGN KEY (user_nickname) REFERENCES users(nickname) ON DELETE CASCADE,
    FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
    FOREIGN KEY (forum_slug) REFERENCES forums(slug) ON DELETE CASCADE
);

CREATE INDEX ON users(lower(nickname));

CREATE
OR REPLACE FUNCTION vote() RETURNS TRIGGER AS $vote$ BEGIN
UPDATE
    threads
SET
    votes = votes + new.voice
WHERE
    id = new.thread_id;
RETURN new;
END;
$vote$ LANGUAGE plpgsql;

CREATE
OR REPLACE FUNCTION revote() RETURNS TRIGGER AS $revote$ BEGIN IF (old.voice != new.voice) THEN
UPDATE
    threads
SET
    votes = votes + new.voice - old.voice
WHERE
    id = new.thread_id;

END IF;
RETURN new;
END;
$revote$ LANGUAGE plpgsql;

CREATE TRIGGER vote_count
AFTER
INSERT
    ON votes FOR EACH ROW EXECUTE PROCEDURE vote();

CREATE TRIGGER revote_count BEFORE
UPDATE
    ON votes FOR EACH ROW EXECUTE PROCEDURE revote();

CREATE
OR REPLACE FUNCTION parent_path() RETURNS trigger LANGUAGE plpgsql AS $func$ BEGIN NEW.path := NEW.path || ARRAY [NEW.id];
RETURN NEW;
END $func$;

CREATE TRIGGER parent_path_count BEFORE
INSERT
    ON posts FOR EACH ROW EXECUTE PROCEDURE parent_path();

CREATE OR REPLACE FUNCTION forum_posts_count()
    RETURNS trigger AS
$forum_posts_count$
BEGIN
    update forums set posts_num = posts_num + 1 where slug = new.forum_slug;
    return new;
END;
$forum_posts_count$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION forum_thread_count()
    RETURNS trigger AS
$forum_thread_count$
BEGIN
    update forums set threads_num = threads_num + 1 where slug = new.forum_slug;
    return new;
END;
$forum_thread_count$ LANGUAGE plpgsql;

CREATE TRIGGER forum_post_count
    AFTER INSERT ON posts
    FOR EACH ROW
EXECUTE PROCEDURE forum_posts_count();

CREATE TRIGGER forum_thread_count
AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE forum_thread_count();

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO db_forum_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO db_forum_user;