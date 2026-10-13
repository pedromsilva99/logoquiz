-- +goose Up
CREATE TABLE levels (
    id     SERIAL PRIMARY KEY,
    number INT  NOT NULL UNIQUE,
    title  TEXT NOT NULL
);

CREATE TABLE logos (
    id        SERIAL PRIMARY KEY,
    level_id  INT  NOT NULL REFERENCES levels (id) ON DELETE CASCADE,
    position  INT  NOT NULL,
    name      TEXT NOT NULL,
    -- Base filename (without extension) under assets/logos/{hidden,full}/. Never sent to the client.
    image_key TEXT NOT NULL UNIQUE,
    UNIQUE (level_id, position)
);

-- Alternative accepted answers (e.g. "CGD" for "Caixa Geral de Depósitos").
CREATE TABLE logo_aliases (
    logo_id INT  NOT NULL REFERENCES logos (id) ON DELETE CASCADE,
    alias   TEXT NOT NULL,
    PRIMARY KEY (logo_id, alias)
);

CREATE TABLE players (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE progress (
    player_id UUID NOT NULL REFERENCES players (id) ON DELETE CASCADE,
    logo_id   INT  NOT NULL REFERENCES logos (id) ON DELETE CASCADE,
    attempts  INT  NOT NULL DEFAULT 0,
    solved_at TIMESTAMPTZ,
    PRIMARY KEY (player_id, logo_id)
);

-- +goose Down
DROP TABLE progress;
DROP TABLE players;
DROP TABLE logo_aliases;
DROP TABLE logos;
DROP TABLE levels;
