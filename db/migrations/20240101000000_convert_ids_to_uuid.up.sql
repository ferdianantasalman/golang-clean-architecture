ALTER TABLE users MODIFY id CHAR(36) NOT NULL;

ALTER TABLE contacts MODIFY id CHAR(36) NOT NULL,
                     MODIFY user_id CHAR(36) NOT NULL;

ALTER TABLE addresses MODIFY id CHAR(36) NOT NULL,
                      MODIFY contact_id CHAR(36) NOT NULL;
