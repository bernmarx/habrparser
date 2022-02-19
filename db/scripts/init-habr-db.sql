CREATE TABLE IF NOT EXISTS habr_pages (
    id              INT PRIMARY KEY,
    title           TEXT NOT NULL,
    article         TEXT NOT NULL,
    posted          TIMESTAMP WITH TIME ZONE,
    author          VARCHAR(64),
    comment_count   INT CHECK (comment_count >= 0),
    rating          INT
);

CREATE OR REPLACE FUNCTION addpage(json)
    RETURNS void
    LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
    INSERT INTO habr_pages
	SELECT *
    FROM json_populate_record(null::habr_pages, $1)
    ON CONFLICT
    DO NOTHING;
END;
$BODY$;