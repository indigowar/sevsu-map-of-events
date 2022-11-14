CREATE TABLE organizer_level
(
    organizer_level_id   UUID PRIMARY KEY,
    organizer_level_name VARCHAR(255),
    organizer_level_code VARCHAR(3)
);

CREATE TABLE organizer
(
    organizer_id    UUID PRIMARY KEY,
    organizer_name  VARCHAR(255) NOT NULL,
    organizer_image VARCHAR(255) NOT NULL,
    organizer_level UUID         NOT NULL,
    FOREIGN KEY (organizer_level) REFERENCES organizer_level (organizer_level_id)
);

CREATE TABLE competitor
(
    competitor_id   UUID PRIMARY KEY,
    competitor_name VARCHAR(255) NOT NULL
);

CREATE TABLE founding_range
(
    founding_range_id   UUID PRIMARY KEY,
    founding_range_low  INT NOT NULL DEFAULT 0,
    founding_range_high INT NOT NULL DEFAULT 0
);

CREATE TABLE co_founding_range
(
    co_founding_range_id UUID PRIMARY KEY,
    co_founding_low      INT NOT NULL DEFAULT 0,
    co_founding_high     INT NOT NULL DEFAULT 0
);

CREATE TABLE event
(
    event_id                   UUID PRIMARY KEY,
    event_organizer            UUID          NOT NULL,
    FOREIGN KEY (event_organizer) REFERENCES organizer (organizer_id),
    event_founding_type        VARCHAR(1024) NOT NULL,
    event_founding_range       UUID          NOT NULL,
    FOREIGN KEY (event_founding_range) REFERENCES founding_range (founding_range_id),
    event_co_founding_range    UUID          NOT NULL,
    FOREIGN KEY (event_co_founding_range) REFERENCES co_founding_range (co_founding_range_id),
    event_submission_deadline  TIME          NOT NULL,
    event_consideration_period VARCHAR(255),
    event_realisation_period   VARCHAR(255),
    event_result               TEXT,
    event_site                 VARCHAR(1024),
    event_document             VARCHAR(1024),
    event_internal_contacts    VARCHAR(255),
    event_trl                  INT           NOT NULL DEFAULT 5
);

CREATE TABLE subject
(
    subject_id    UUID PRIMARY KEY,
    subject_name  VARCHAR(255) NOT NULL,
    subject_event UUID,
    FOREIGN KEY (subject_event) REFERENCES event (event_id)
);

CREATE TABLE competitor_requirements
(
    cr_id         UUID PRIMARY KEY,
    cr_event      UUID,
    FOREIGN KEY (cr_event) REFERENCES event (event_id),
    cr_competitor UUID,
    FOREIGN KEY (cr_competitor) REFERENCES competitor (competitor_id)
)
