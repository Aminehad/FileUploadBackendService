CREATE DATABASE file_upload;
CREATE USER user_file_upload WITH PASSWORD 'pass';

ALTER DATABASE file_upload OWNER TO user_file_upload;
GRANT ALL PRIVILEGES ON DATABASE file_upload TO user_file_upload;
GRANT CONNECT ON DATABASE file_upload TO user_file_upload;
GRANT USAGE, CREATE ON SCHEMA public TO user_file_upload;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO user_file_upload;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO user_file_upload;
