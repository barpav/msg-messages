CREATE SEQUENCE timeline AS bigint; 

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    event_timestamp bigint DEFAULT nextval('timeline') UNIQUE NOT NULL,
    sender varchar(50) NOT NULL,
    receiver varchar(50) NOT NULL,
    created timestamp,
    edited timestamp,
    read_at timestamp,
    message_text text,
    is_deleted bool
);

CREATE TABLE attachments (
    message_id bigint REFERENCES messages(id) NOT NULL,
    file_id varchar(24) NOT NULL
);

CREATE INDEX attachments_idx ON attachments (message_id);

CREATE TABLE updates (
    user_id varchar(50) NOT NULL,
    event_timestamp bigint NOT NULL,
    message_id bigint REFERENCES messages(id) NOT NULL
);

CREATE UNIQUE INDEX updates_idx ON updates (user_id, event_timestamp);