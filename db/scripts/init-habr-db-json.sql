CREATE TABLE IF NOT EXISTS habr_pages2 (
    id          INT PRIMARY KEY,
    pageinfo    JSONB
);

CREATE OR REPLACE FUNCTION addpagejson(INT, JSONB)
    RETURNS void
    LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
    INSERT INTO habr_pages2
	VALUES (
        $1,
        to_jsonb($2)
    )
    ON CONFLICT
    DO NOTHING;
END;
$BODY$;