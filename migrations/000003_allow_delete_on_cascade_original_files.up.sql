ALTER TABLE translated_files DROP CONSTRAINT translated_files_original_files_id_fkey;

ALTER TABLE translated_files
        ADD CONSTRAINT translated_files_original_files_id_fkey
        FOREIGN KEY (original_files_id) REFERENCES original_files(id) ON DELETE CASCADE;
