-- name: GetProductByID :one
SELECT id, account_id, name
FROM products
WHERE account_id = $1 AND id = $2 AND deleted_at IS NULL;

-- name: ListProducts :many
SELECT id, name FROM products WHERE account_id = $1;

-- name: ListProductsPage :paginated
SELECT id, name FROM products WHERE account_id = $1;

-- name: ListProductsByName :paginated
SELECT id, name FROM products WHERE account_id = $1 AND deleted_at IS NULL ORDER BY name;

-- name: ListProductsIncludingDeleted :many
SELECT id, name FROM products WHERE account_id = $1;

-- name: GetProductAudit :many
SELECT id, name, deleted_at FROM products WHERE account_id = $1;

-- name: SoftDeleteProduct :exec
UPDATE products SET deleted_at = now() WHERE id = $1;

-- name: InsertProduct :one
INSERT INTO products (id, name) VALUES ($1, $2) RETURNING id;
