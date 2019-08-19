-- ----------------------------------------------------------------------------
-- Users
-- ----------------------------------------------------------------------------

CREATE TABLE Users (
    id bigserial PRIMARY KEY,
    username text NOT NULL UNIQUE,
    password_hash text NOT NULL,

    usertype smallint NOT NULL,
    signature text,
    about text,
    banned boolean,

    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
ALTER TABLE Users OWNER TO riftforum_user;

-- ----------------------------------------------------------------------------
-- Invites
-- ----------------------------------------------------------------------------

CREATE TABLE Invites (
    id bigserial PRIMARY KEY,
    key text NOT NULL UNIQUE,
    status smallint NOT NULL,

    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
ALTER TABLE Invites OWNER TO riftforum_user;

-- ----------------------------------------------------------------------------
-- Topics
-- ----------------------------------------------------------------------------

CREATE TABLE Topics (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    author_id bigserial NOT NULL,

    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),

    FOREIGN KEY (author_id) REFERENCES Users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
ALTER TABLE Topics OWNER TO riftforum_user;

-- ----------------------------------------------------------------------------
-- Messages
-- ----------------------------------------------------------------------------

CREATE TABLE Messages (
    id bigserial PRIMARY KEY,
    author_id bigserial NOT NULL,
    topic_id bigserial NOT NULL,
    message text NOT NULL,

    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),

    FOREIGN KEY (author_id) REFERENCES Users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (topic_id) REFERENCES Topics(id) ON DELETE CASCADE ON UPDATE CASCADE
);
ALTER TABLE Messages OWNER TO riftforum_user;
