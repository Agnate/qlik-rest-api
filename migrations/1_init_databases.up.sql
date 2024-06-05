CREATE TABLE IF NOT EXISTS users
(
    uuid uuid NOT NULL DEFAULT gen_random_uuid(),
    email character varying(350) COLLATE pg_catalog."default",
    api_key text COLLATE pg_catalog."default",
    last_access timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    create_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    full_name character varying(150) COLLATE pg_catalog."default",
    CONSTRAINT users_pkey PRIMARY KEY (uuid),
    CONSTRAINT unique_api_key UNIQUE NULLS NOT DISTINCT (api_key),
    CONSTRAINT unique_email UNIQUE NULLS NOT DISTINCT (email)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS users
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS messages
(
    uuid uuid NOT NULL,
    create_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    message text COLLATE pg_catalog."default",
    is_palindrome boolean,
    last_updated timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    last_updated_by uuid,
    logical_delete boolean DEFAULT false,
    CONSTRAINT messages_pkey PRIMARY KEY (uuid, create_date),
    CONSTRAINT last_updated_by FOREIGN KEY (last_updated_by)
        REFERENCES users (uuid) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT uuid FOREIGN KEY (uuid)
        REFERENCES users (uuid) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS messages
    OWNER to postgres;

INSERT INTO users (email, api_key, full_name) VALUES ('bobross@gmail.com', 'very-secure-api-key', 'Bob Ross');
INSERT INTO users (email, api_key, full_name) VALUES ('mrrogers@gmail.com', 'another-equally-secure-api-key', 'Mr Rogers');
