use anti_brute_force;


CREATE TABLE IF NOT EXISTS masks (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  list_id SMALLINT NOT NULL,
  mask CHAR(18)  NOT NULL,
  INDEX ix_masks_list_id (list_id),
  UNIQUE INDEX ux_masks_list_id_mask (list_id, mask)
);