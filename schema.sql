\c forum;

DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE users (
    id serial,
    nickname citext COLLATE "C" PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    fullname text NOT NULL,
    about text
);

CREATE UNLOGGED  TABLE forums (
    slug citext PRIMARY KEY,
    title text NOT NULL,
    user_nickname citext NOT NULL,
    threads_num int NOT NULL DEFAULT 0,
    posts_num int NOT NULL DEFAULT 0,
    FOREIGN KEY (user_nickname) REFERENCES users(nickname) ON DELETE CASCADE
);

CREATE UNLOGGED  TABLE threads (
    id serial PRIMARY KEY,
    slug citext NOT NULL,
    title text NOT NULL,
    message text NOT NULL,
    created timestamptz NOT NULL DEFAULT NOW(),
    user_nickname citext NOT NULL,
    forum_slug citext NOT NULL,
    votes int NOT NULL DEFAULT 0,
    FOREIGN KEY (user_nickname) REFERENCES users(nickname) ON DELETE CASCADE,
    FOREIGN KEY (forum_slug) REFERENCES forums(slug) ON DELETE CASCADE
);

CREATE UNLOGGED  TABLE votes (
    user_nickname citext NOT NULL REFERENCES users (nickname) ON DELETE CASCADE,
    thread_id int NOT NULL REFERENCES threads (id) ON DELETE CASCADE,
    voice int NOT NULL,
    PRIMARY KEY(user_nickname, thread_id)
);

CREATE UNLOGGED  TABLE posts (
    id serial PRIMARY KEY,
    message text NOT NULL,
    is_edited boolean NOT NULL DEFAULT false,
    created timestamptz NOT NULL DEFAULT NOW(),
    parent_id int NOT NULL DEFAULT 0,
    path integer [] NOT NULL DEFAULT '{}',
    user_nickname citext NOT NULL,
    forum_slug citext NOT NULL,
    thread_id int NOT NULL,
    FOREIGN KEY (user_nickname) REFERENCES users(nickname) ON DELETE CASCADE,
    FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
    FOREIGN KEY (forum_slug) REFERENCES forums(slug) ON DELETE CASCADE
);

CREATE INDEX ON posts(thread_id);
CREATE INDEX ON posts(path);
CREATE INDEX ON posts(user_nickname);
CREATE INDEX ON posts(forum_slug);
CREATE INDEX ON posts(id, thread_id);

CREATE INDEX ON forums(user_nickname);

CREATE INDEX ON votes(user_nickname);
CREATE INDEX ON votes(thread_id);

CREATE INDEX ON threads(user_nickname);
CREATE INDEX ON threads(forum_slug);
CREATE INDEX ON threads(id);
CREATE INDEX ON threads(slug);
CREATE INDEX ON threads(id, slug, title, message, created, user_nickname, forum_slug, votes) ;

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

CREATE FUNCTION trigger_post_before_insert()
    RETURNS trigger AS $trigger_post_before_insert$
BEGIN
    IF NEW.path[1] <> 0 THEN
        NEW.path := (SELECT path FROM posts WHERE id = NEW.path[1] AND thread_id = NEW.thread_id) || ARRAY[NEW.id];
        IF array_length(NEW.path, 1) = 1 THEN
            RAISE 'Parent post exception' USING ERRCODE = '12345';
        END IF;
    ELSE
        NEW.path[1] := NEW.id;
    END IF;
    RETURN NEW;
END;
$trigger_post_before_insert$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE trigger_post_before_insert();


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
