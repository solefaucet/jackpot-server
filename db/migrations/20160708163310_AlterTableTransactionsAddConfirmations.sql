
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `transactions` ADD COLUMN `confirmations` INT(11) NOT NULL COMMENT 'number of confirmations of the tranactions'; 
ALTER TABLE `transactions` ADD INDEX (`confirmations`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `transactions` DROP COLUMN `confirmations`;
