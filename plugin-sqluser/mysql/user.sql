CREATE TABLE user(
    user_id VARCHAR(255) not null,
    created_time BIGINT not null,
    banned int not null,
    PRIMARY KEY(id)
) DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci  ENGINE=InnoDB; 