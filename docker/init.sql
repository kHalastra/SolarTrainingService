-- Create the 'books' table only if it does not already exist
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    year INT NOT NULL
);

-- Insert data into the 'books' table, but only if it's not already inserted
-- You can use ON CONFLICT to avoid duplication if necessary
INSERT INTO books (title, author, year)
SELECT 'Lord of the rings', 'J.R.R. Tolkien', 2004
WHERE NOT EXISTS (SELECT 1 FROM books WHERE title = 'Lord of the rings');

INSERT INTO books (title, author, year)
SELECT 'Hounded', 'Kevin Hearne', 2010
WHERE NOT EXISTS (SELECT 1 FROM books WHERE title = 'Hounded');
