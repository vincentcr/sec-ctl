BEGIN;

DROP FUNCTION IF EXISTS normalize_uuid(uuid);

CREATE FUNCTION normalize_uuid(uuid) RETURNS TEXT AS
$$
  SELECT replace(CAST($1 AS TEXT), '-', '')
$$ LANGUAGE sql;


COMMIT;
