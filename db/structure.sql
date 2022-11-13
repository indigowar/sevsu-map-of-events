CREATE TABLE organizer_level
(
    id   UUID PRIMARY KEY,
    name VARCHAR(255),
    code VARCHAR(3)
);

CREATE TABLE organizer
(
    id    UUID PRIMARY KEY,
    name  VARCHAR(255)  NOT NULL,
    logo  VARCHAR(1024) NOT NULL,
    level UUID          NOT NULL,
    FOREIGN KEY (level) REFERENCES organizer_level (id)
);

CREATE TABLE subject
(
    id   UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE competitor
(
    id   UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE founding_range
(
    id   UUID PRIMARY KEY,
    low  INT NOT NULL DEFAULT 0,
    high INT NOT NULL DEFAULT 0
);

CREATE TABLE co_founding_range
(
    id   UUID PRIMARY KEY,
    low  INT NOT NULL DEFAULT 0,
    high INT NOT NULL DEFAULT 0
);

CREATE TABLE event
(
    id                   UUID PRIMARY KEY,
    organizer            UUID          NOT NULL,
    founding_type        VARCHAR(1024) NOT NULL,
    founding_range       UUID          NOT NULL,
    co_founding_range    UUID          NOT NULL,
    submission_deadline  TIME          NOT NULL,
    consideration_period VARCHAR(255),
    realisation_period   VARCHAR(255),
    result               TEXT,
    site                 VARCHAR(1024),
    document             VARCHAR(1024),
    internal_contacts    VARCHAR(255),
    trl_value            INT           NOT NULL DEFAULT 5
);

CREATE TABLE subject
(
    id    UUID PRIMARY KEY,
    name  VARCHAR(255) NOT NULL,
    event UUID,
    FOREIGN KEY (event) REFERENCES event (id)
);

CREATE TABLE competitor_requirements
(
    id         UUID PRIMARY KEY,
    event      UUID,
    FOREIGN KEY (event) REFERENCES event (id),

    competitor UUID,
    FOREIGN KEY (competitor) REFERENCES competitor (id)
)
