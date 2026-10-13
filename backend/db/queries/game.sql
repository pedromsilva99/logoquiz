-- name: CreatePlayer :one
INSERT INTO players DEFAULT VALUES
RETURNING id;

-- name: PlayerExists :one
SELECT EXISTS (SELECT 1 FROM players WHERE id = $1);

-- name: GetLevelByNumber :one
SELECT id, number, title
FROM levels
WHERE number = $1;

-- name: ListLevelLogosForPlayer :many
SELECT lg.id, lg.position, lg.name, (p.solved_at IS NOT NULL)::bool AS solved
FROM logos lg
LEFT JOIN progress p ON p.logo_id = lg.id AND p.player_id = sqlc.arg(player_id)
WHERE lg.level_id = sqlc.arg(level_id)
ORDER BY lg.position;

-- name: GetLogo :one
SELECT id, name, image_key
FROM logos
WHERE id = $1;

-- name: ListLogoAliases :many
SELECT alias
FROM logo_aliases
WHERE logo_id = $1;

-- name: IsLogoSolved :one
SELECT EXISTS (
    SELECT 1 FROM progress
    WHERE player_id = $1 AND logo_id = $2 AND solved_at IS NOT NULL
);

-- name: RecordAttempt :exec
INSERT INTO progress (player_id, logo_id, attempts, solved_at)
VALUES (
    sqlc.arg(player_id),
    sqlc.arg(logo_id),
    1,
    CASE WHEN sqlc.arg(correct)::bool THEN now() END
)
ON CONFLICT (player_id, logo_id) DO UPDATE
SET attempts  = progress.attempts + 1,
    solved_at = COALESCE(progress.solved_at, EXCLUDED.solved_at);
