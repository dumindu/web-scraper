-- +goose Up

CREATE TABLE "users"
(
    "id"         UUID                     NOT NULL,
    "email"      TEXT UNIQUE              NOT NULL,
    "created_at" TIMESTAMP with time zone NOT NULL,
    "updated_at" TIMESTAMP with time zone NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "user_auths"
(
    "user_id"    UUID                     NOT NULL,
    "password"   TEXT                     NOT NULL,
    "created_at" TIMESTAMP with time zone NOT NULL,
    "updated_at" TIMESTAMP with time zone NOT NULL,
    PRIMARY KEY ("user_id")
);

CREATE TABLE "user_activation_tokens"
(
    "user_id"          UUID                     NOT NULL,
    "password"         TEXT                     NOT NULL,
    "token_expired_at" TIMESTAMP with time zone NOT NULL,
    PRIMARY KEY ("user_id")
);

CREATE TABLE "user_reset_password_tokens"
(
    "user_id"          UUID                     NOT NULL,
    "password"         TEXT                     NOT NULL,
    "token_expired_at" TIMESTAMP with time zone NOT NULL,
    PRIMARY KEY ("user_id")
);

CREATE TABLE "keywords"
(
    "id"            SERIAL                   NOT NULL,
    "user_id"       UUID                     NOT NULL,
    "keyword"       TEXT                     NOT NULL,
    "status"        TEXT                     NOT NULL DEFAULT 'pending',
    "search_engine" TEXT                     NOT NULL,
    "ad_count"      BIGINT,
    "link_count"    BIGINT,
    "html_content"  TEXT,
    "error_message" TEXT,
    "created_at"    TIMESTAMP with time zone NOT NULL,
    "updated_at"    TIMESTAMP with time zone NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT uix_keyword_user_keyword_search_engine UNIQUE ("user_id", "keyword", "search_engine")
);

-- +goose Down

DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "user_auths";
DROP TABLE IF EXISTS "user_activation_tokens";
DROP TABLE IF EXISTS "user_reset_password_tokens";
DROP TABLE IF EXISTS "keywords";
