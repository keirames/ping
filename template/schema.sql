CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

CREATE TABLE users (
	id BIGINT PRIMARY KEY,
	name VARCHAR(255) NOT NULL
);

CREATE TABLE chat_rooms (
	id BIGINT PRIMARY KEY,
	name varchar(255) NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW()
);


CREATE TABLE messages (
	id BIGINT PRIMARY KEY,
	content TEXT NOT NULL,
	type TEXT CHECK (TYPE IN ('image', 'text', 'emote')) DEFAULT 'text',
	is_delete BOOLEAN DEFAULT false,
	parent_id BIGINT,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	
	user_id BIGINT NOT NULL REFERENCES users(id),
	room_id BIGINT NOT NULL REFERENCES chat_rooms(id),
	
	CHECK (content <> '')
);

CREATE TABLE users_and_chat_rooms (
	id BIGINT PRIMARY KEY,
	
	user_id BIGINT NOT NULL REFERENCES users(id),
	room_id BIGINT NOT NULL REFERENCES chat_rooms(id),

	UNIQUE (user_id, room_id)
);

CREATE TABLE message_emote_types (
	id BIGINT PRIMARY KEY,
	
	name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE messages_emotes (
	id BIGINT PRIMARY KEY,
	
	message_id BIGINT NOT NULL REFERENCES messages(id),
	user_id BIGINT NOT NULL REFERENCES users(id),
	emote_type_id BIGINT NOT NULL REFERENCES message_emote_types(id),
	
	UNIQUE (message_id, user_id, emote_type_id)
);