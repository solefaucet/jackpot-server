
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `games` ADD COLUMN `status` VARCHAR(255) NOT NULL DEFAULT 'pending' COMMENT 'game status';
ALTER TABLE `games` ADD INDEX (`status`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `games` DROP COLUMN `status`;
