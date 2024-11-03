CREATE TABLE IF NOT EXISTS feed_follows(
   id UUID PRIMARY KEY,
   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP NOT NULL,

   feed_id UUID REFERENCES feeds ON DELETE CASCADE,
   user_id UUID REFERENCES users ON DELETE CASCADE,
   UNIQUE (user_id, feed_id)
);
