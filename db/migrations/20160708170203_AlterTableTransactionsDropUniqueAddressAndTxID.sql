
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `transactions` DROP INDEX `address`;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `transactions` ADD UNIQUE INDEX (`address`, `tx_id`);
