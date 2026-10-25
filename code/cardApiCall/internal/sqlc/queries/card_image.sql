-- name: InsertCardImage :one
INSERT INTO card_images (
    ocg_api_image_id,
    ocg_api_card_id,
    card_id,
    image_url_small_api,
    image_url_large_api,
    image_url_cropped_api
) VALUES (
    $1, $2, $3, $4, $5, $6
)
ON CONFLICT (ocg_api_image_id) DO UPDATE
SET
    card_id = EXCLUDED.card_id,
    image_url_small_api = EXCLUDED.image_url_small_api,
    image_url_large_api = EXCLUDED.image_url_large_api,
    image_url_cropped_api = EXCLUDED.image_url_cropped_api
RETURNING *;

-- name: ListCardImagesByCardID :many
SELECT * FROM card_images
WHERE card_id = $1
ORDER BY id;

-- name: ListCardImagesByOcgApiCardID :many
SELECT ci.* FROM card_images ci
LEFT JOIN cards c ON ci.card_id = c.id
WHERE ci.ocg_api_card_id = $1 OR c.neuron_id = $1
ORDER BY ci.id;

-- name: ListPendingSyncCardImages :many
SELECT * FROM card_images
WHERE image_url_large_this_app IS NULL
ORDER BY id
LIMIT $1;

-- name: UpdateCardImageThisAppUrls :one
UPDATE card_images
SET
    image_url_small_this_app = $2,
    image_url_large_this_app = $3,
    image_url_cropped_this_app = $4
WHERE id = $1
RETURNING *;

-- name: DeleteCardImagesByCardID :exec
DELETE FROM card_images
WHERE card_id = $1;
