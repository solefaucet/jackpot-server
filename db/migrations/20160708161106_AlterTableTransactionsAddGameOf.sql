
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `transactions` ADD COLUMN `game_of` DATETIME NOT NULL COMMENT 'game of time'; 
ALTER TABLE `transactions` ADD INDEX (`game_of`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `transactions` DROP COLUMN `game_of`;
