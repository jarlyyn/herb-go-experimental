CREATE TABLE account(
    uid VARCHAR(255) not null,
    keyword VARCHAR(255) not null,
    account VARCHAR(255)
    CHARACTER SET utf8 
    COLLATE utf8_bin
    not null,
    created_time BIGINT not null,
    PRIMARY KEY(keyword,account)
  
) DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ENGINE=InnoDB;