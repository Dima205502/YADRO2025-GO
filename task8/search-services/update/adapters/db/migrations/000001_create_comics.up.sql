CREATE TABLE comics (
    id SERIAL PRIMARY KEY,
    comics_id INTEGER UNIQUE NOT NULL,
    img_url TEXT NOT NULL,
    keywords TEXT[]
);

