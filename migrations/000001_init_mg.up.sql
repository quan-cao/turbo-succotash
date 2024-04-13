CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        isid VARCHAR(255) NOT NULL,
        role VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        UNIQUE (isid)
);

CREATE TABLE IF NOT EXISTS original_files (
        id SERIAL PRIMARY KEY,
        sha256 char(64) NOT NULL,
        filename VARCHAR(255) NOT NULL,
        file_type VARCHAR(255) NOT NULL,
        file_size BIGINT NOT NULL,
        source_language VARCHAR(255) NOT NULL,
        token_count INT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        created_by INT NOT NULL,
        FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS translated_files (
        id SERIAL PRIMARY KEY,
        original_files_id INT,
        translated_filename VARCHAR(255) NOT NULL,
        target_language VARCHAR(255) NOT NULL,
        cost DECIMAL(10, 2) NOT NULL,
        time_taken INT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        created_by INT NOT NULL,
        FOREIGN KEY (original_files_id) REFERENCES original_files(id),
        FOREIGN KEY (created_by) REFERENCES users(id)
);
