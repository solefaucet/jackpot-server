
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `transactions` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `address` VARCHAR(255) NOT NULL COMMENT 'address',
  `amount` DECIMAL(19, 8) NOT NULL COMMENT 'transaction amount',
  `tx_id` VARCHAR(255) NOT NULL COMMENT 'transaction id',
  `hash` VARCHAR(255) NOT NULL COMMENT 'block hash',
  `block_created_at` DATETIME NOT NULL COMMENT 'block created at',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `transactions`
ADD UNIQUE INDEX (`address`, `tx_id`),
ADD INDEX (`address`),
ADD INDEX (`tx_id`),
ADD INDEX (`hash`),
ADD INDEX (`block_created_at`),
ADD INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `transactions`;
