
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `games` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `hash` VARCHAR(255) NOT NULL COMMENT 'block hash',
  `height` INT(11) NOT NULL COMMENT 'block height',
  `address` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'winner address, can be empty',
  `win_amount` DECIMAL(19, 8) NOT NULL DEFAULT 0 COMMENT 'win amount',
  `total_amount` DECIMAL(19, 8) NOT NULL COMMENT 'total amount',
  `fee` DECIMAL(19, 8) NOT NULL DEFAULT 0 COMMENT 'transaction fee',
  `tx_id` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'transaction proof',
  `game_of` DATETIME NOT NULL COMMENT 'game of time',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `games`
ADD UNIQUE INDEX (`game_of`),
ADD UNIQUE INDEX (`height`, `hash`, `game_of`),
ADD INDEX (`height`),
ADD INDEX (`hash`),
ADD INDEX (`created_at`),
ADD INDEX (`address`),
ADD INDEX (`tx_id`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `games`;
