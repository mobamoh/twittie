CREATE TABLE IF NOT EXISTS tweets
(
    id         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v1(),
    body       VARCHAR(250)     NOT NULL,
    user_Id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE  ,
    created_at TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);