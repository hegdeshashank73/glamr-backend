CREATE TABLE auth_users (
  id bigint NOT NULL,
    email varchar(256) NOT NULL,
      PRIMARY KEY (id),
        CONSTRAINT auth_users_email_unq_idx UNIQUE (email)
        );

CREATE TABLE auth_tokens (
    token char(36) NOT NULL,
    user_id bigint,
    created_at bigint NOT NULL,
    PRIMARY KEY (token),
    CONSTRAINT auth_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES auth_users (id)
);

CREATE INDEX auth_tokens_user_id_idx ON auth_tokens (user_id);

CREATE TABLE auth_magiclink (
    token char(36) NOT NULL,
    email varchar(256) NOT NULL,
    created_at bigint NOT NULL,
    PRIMARY KEY (token)
);

CREATE TABLE people_people (
    id bigint NOT NULL,
    first_name varchar(64) NOT NULL DEFAULT '',
    last_name varchar(64) NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    CONSTRAINT people_people_id_fkey FOREIGN KEY (id) 
        REFERENCES auth_users (id) ON DELETE CASCADE
);

CREATE TABLE people_searches (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    s3_key varchar(128) NOT NULL DEFAULT '',
    created_at bigint NOT NULL,
    country_code varchar(16) NOT NULL DEFAULT '',
    api_response text NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT people_searches_people_people_user_id_fkey 
        FOREIGN KEY (user_id) REFERENCES people_people (id) ON DELETE CASCADE
);

CREATE INDEX people_searches_user_id_idx ON people_searches (user_id);

CREATE TABLE searches_options (
    id bigint NOT NULL,
    search_id bigint NOT NULL,
    title varchar(512) NOT NULL DEFAULT '',
    link varchar(512) NOT NULL DEFAULT '',
    source varchar(64) NOT NULL DEFAULT '',
    source_icon varchar(512) NOT NULL DEFAULT '',
    in_stock smallint NOT NULL DEFAULT 0,
    price integer NOT NULL DEFAULT 0,
    currency varchar(8) NOT NULL DEFAULT '',
    rank integer NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    CONSTRAINT searches_options_search_id_fkey 
        FOREIGN KEY (search_id) REFERENCES people_searches (id) ON DELETE CASCADE
);

CREATE INDEX searches_options_search_id_idx ON searches_options (search_id);

CREATE TABLE templates_emails (
    id bigint NOT NULL,
    name varchar(64) NOT NULL,
    body text NOT NULL,
    subject text NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT unq_templates_emails_name UNIQUE (name)
);

ALTER TABLE searches_options RENAME COLUMN rank TO display_order;
ALTER TABLE searches_options ADD COLUMN image varchar(256) NOT NULL DEFAULT '';

INSERT INTO auth_users (id,email) values (243725088341860353, 'hegdeshashank73@gmail.com');
INSERT INTO people_people (id) values (243725088341860353); 
INSERT INTO auth_tokens(user_id, token, created_at) 
VALUES (243725088341860353, '5e96a75b-c194-4f6e-9f26-3ccd509c9755', 1739693279);