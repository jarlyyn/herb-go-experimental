CREATE TABLE revoke_token(
    uid VARCHAR(255) not null,
    created_time BIGINT not null,
    revoke_token VARCHAR(255),
    PRIMARY KEY(id)
) DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ENGINE=InnoDB; 