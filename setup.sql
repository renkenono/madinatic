USE madinatic;


-- drop dependent tables first
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS upvotes;
DROP TABLE IF EXISTS subs;
DROP TABLE IF EXISTS repcats;
DROP TABLE IF EXISTS bans;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS pictures;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS citizens;
DROP TABLE IF EXISTS authorities;
DROP TABLE IF EXISTS users;


CREATE TABLE users (
    pk_userid BIGINT,
    username VARCHAR(30) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,

    -- 213 x xx xx xx xx
    -- devs must add the "+" sign
    -- int provides faster comparision and search
    -- The length just specifies how many characters
    -- to display when selecting data with the mysql command line client.
    phone BIGINT NOT NULL UNIQUE,

    -- it seems that bcrypt always generates 60 character hashes.
    -- bcrypt limits your max password length to 50-72 bytes
    -- depending on the implementation, thus 50 for safety
    password VARCHAR(60) NOT NULL,

    -- ref: http://go-database-sql.org/nulls.html
    -- NULL should be avoided when possible,
    -- empty string "" will be used in this matter
    confirm_token VARCHAR(255) NOT NULL,
    reset_token VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT NOW(),

    -- devs will have to expilicity check if
    -- modified_at is equal to created_at. If so,
    -- no new modifications were introduced
    modified_at DATETIME DEFAULT NOW(),
    PRIMARY KEY (pk_userid)
);

CREATE TABLE citizens (
    pk_userid BIGINT,
    first_name VARCHAR(30) NOT NULL,
    family_name VARCHAR(30) NOT NULL,
    PRIMARY KEY (pk_userid),
    FOREIGN KEY (pk_userid)
    REFERENCES users(pk_userid)
    ON DELETE CASCADE
);

CREATE TABLE authorities (
    pk_userid BIGINT,
    name VARCHAR(30) NOT NULL,
    PRIMARY KEY (pk_userid),
    FOREIGN KEY (pk_userid)
    REFERENCES users(pk_userid)
    ON DELETE CASCADE
);

CREATE TABLE reports (
    pk_reportid INT AUTO_INCREMENT,
    title VARCHAR(80) NOT NULL,
    descr TEXT NOT NULL,
    created_at DATETIME DEFAULT NOW(),
    modified_at DATETIME DEFAULT NOW(),
    latitude DOUBLE NOT NULL,
    longitude DOUBLE NOT NULL,

    -- waiting, rejected, processing, solved [0, 1, 2, 3]
    curr_state TINYINT CHECK (curr_state >= 0 AND curr_state < 4),

    -- make sure to compare it to created_at
    state_modified_at DATETIME DEFAULT NOW(),
    fk_userid BIGINT,
    PRIMARY KEY (pk_reportid),
    FOREIGN KEY (fk_userid) REFERENCES users(pk_userid)
);

CREATE TABLE pictures (

    -- low chance that names would be compared, string is fine
    pk_name VARCHAR(20),
    fk_reportid INT,
    PRIMARY KEY (pk_name),
    FOREIGN KEY (fk_reportid) REFERENCES reports(pk_reportid)
);

CREATE TABLE categories (
    pk_catid INT AUTO_INCREMENT,
    cat_name VARCHAR(80) NOT NULL,
    fk_userid BIGINT,
    PRIMARY KEY (pk_catid),
    FOREIGN KEY (fk_userid) REFERENCES users(pk_userid)
);

CREATE TABLE comments (
    fk_reportid INT,
    fk_userid BIGINT,
    comment VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT NOW(),
    PRIMARY KEY (fk_reportid, fk_userid),
    FOREIGN KEY (fk_reportid) REFERENCES reports(pk_reportid),
    FOREIGN KEY (fk_userid) REFERENCES users(pk_userid)
);

CREATE TABLE upvotes (
    fk_reportid INT,
    fk_userid BIGINT,
    PRIMARY KEY (fk_reportid, fk_userid),
    FOREIGN KEY (fk_reportid) REFERENCES reports(pk_reportid),
    FOREIGN KEY (fk_userid) REFERENCES users(pk_userid)
);

CREATE TABLE subs (
    fk_reportid INT,
    fk_userid BIGINT,
    PRIMARY KEY (fk_reportid, fk_userid),
    FOREIGN KEY (fk_reportid) REFERENCES reports(pk_reportid),
    FOREIGN KEY (fk_userid) REFERENCES users(pk_userid)
);

CREATE TABLE repcats (
    fk_reportid INT,
    fk_userid BIGINT,

    -- 0 is false, otherwise true
    solved BOOLEAN DEFAULT 0 NOT NULL,
    PRIMARY KEY (fk_reportid, fk_userid),
    FOREIGN KEY (fk_reportid) REFERENCES reports(pk_reportid),
    FOREIGN KEY (fk_userid) REFERENCES users(pk_userid)
);

CREATE TABLE bans (
    fk_userid BIGINT,
    solved BOOLEAN,
    PRIMARY KEY (fk_userid),
    FOREIGN KEY (fk_userid) REFERENCES users(pk_userid)
);

CREATE TABLE notifications (
    fk_reportid INT,
    fk_userid BIGINT,

    -- type comment, solved, accepted, rejected (might change)
    n_type TINYINT NOT NULL CHECK (n_type >= 0 AND n_type < 4),
    n_message VARCHAR(255) NOT NULL,
    PRIMARY KEY (fk_reportid, fk_userid),
    FOREIGN KEY (fk_reportid) REFERENCES reports(pk_reportid),
    FOREIGN KEY (fk_userid) REFERENCES users(pk_userid)
);
