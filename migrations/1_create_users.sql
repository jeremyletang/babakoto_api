USE babakoto;

CREATE TABLE IF NOT EXISTS users
(
  id                  VARCHAR(36)                                                     NOT NULL,
  username            VARCHAR(512)                                                    NOT NULL,
  email               VARCHAR(512)                                                    NOT NULL,
  password            VARCHAR(512)                                                    NOT NULL,
  created_at          DATETIME DEFAULT CURRENT_TIMESTAMP                              NOT NULL,
  updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS access_tokens
(
  id         VARCHAR(36)                        NOT NULL,
  user_id    VARCHAR(36)                        NOT NULL,
  ttl        INTEGER                            NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS user_signup_verifications
(
  id         VARCHAR(36)                        NOT NULL,
  user_id    VARCHAR(36)                        NOT NULL,
  ttl        INTEGER                            NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE access_tokens
      ADD FOREIGN KEY (user_id) REFERENCES users (id);


