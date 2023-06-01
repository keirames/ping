CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

drop table users cascade;
drop table messages cascade;
drop table chat_rooms cascade;
drop table users_and_chat_rooms cascade;

CREATE TABLE users (
	id BIGINT,
	name VARCHAR(255) NOT NULL,
	
	PRIMARY KEY (id)
);

CREATE TABLE chat_rooms (
	id BIGINT,
	name varchar(255) NOT NULL,
	
	PRIMARY KEY (id)
);


CREATE TABLE messages (
	id BIGINT,
	content TEXT NOT NULL,
	type TEXT CHECK (TYPE IN ('image', 'text', 'emote')) DEFAULT 'text',
	
	user_id BIGINT NOT NULL REFERENCES users(id),
	room_id BIGINT NOT NULL REFERENCES chat_rooms(id),
	
	PRIMARY KEY (id),
	UNIQUE (user_id, room_id)
);

CREATE TABLE users_and_chat_rooms (
	id BIGINT,
	
	user_id BIGINT NOT NULL REFERENCES users(id),
	room_id BIGINT NOT NULL REFERENCES chat_rooms(id),
	
	PRIMARY KEY (id),
	UNIQUE (user_id, room_id)
);