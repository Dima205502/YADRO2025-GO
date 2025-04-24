CREATE TABLE comics (
    id SERIAL PRIMARY KEY,
    comic_id INTEGER UNIQUE NOT NULL,
    img_url TEXT NOT NULL,
    keywords TEXT[]
);

