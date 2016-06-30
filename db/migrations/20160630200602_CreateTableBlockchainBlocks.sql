
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `blocks` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `hash` VARCHAR(255) NOT NULL COMMENT 'block hash',
  `height` INT(11) NOT NULL COMMENT 'block height',
  `block_created_at` DATETIME NOT NULL COMMENT 'block created at',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `blocks`
ADD UNIQUE INDEX (`hash`),
ADD UNIQUE INDEX (`height`),
ADD INDEX (`block_created_at`),
ADD INDEX (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `blocks`;
