ALTER TABLE original_files DROP CONSTRAINT original_files_created_by_fkey;
ALTER TABLE translated_files DROP CONSTRAINT translated_files_created_by_fkey;

ALTER TABLE users ADD CONSTRAINT unique_isid UNIQUE (isid);

ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD PRIMARY KEY (isid);

ALTER TABLE original_files ADD COLUMN new_created_by VARCHAR(255);
UPDATE original_files SET new_created_by = (SELECT isid FROM users WHERE id = original_files.created_by);
ALTER TABLE original_files DROP COLUMN created_by;
ALTER TABLE original_files RENAME COLUMN new_created_by TO created_by;
ALTER TABLE original_files ADD CONSTRAINT fk_original_files_created_by FOREIGN KEY (created_by) REFERENCES users (isid);

ALTER TABLE translated_files ADD COLUMN new_created_by VARCHAR(255);
UPDATE translated_files SET new_created_by = (SELECT isid FROM users WHERE id = translated_files.created_by);
ALTER TABLE translated_files DROP COLUMN created_by;
ALTER TABLE translated_files RENAME COLUMN new_created_by TO created_by;
ALTER TABLE translated_files ADD CONSTRAINT fk_translated_files_created_by FOREIGN KEY (created_by) REFERENCES users (isid);
