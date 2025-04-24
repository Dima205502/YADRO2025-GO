CREATE OR REPLACE FUNCTION array_intersect_count(arr1 text[], arr2 text[])
RETURNS integer AS $$
  SELECT count(*)
  FROM (
    SELECT unnest(arr1)
    INTERSECT
    SELECT unnest(arr2)
  ) AS common;
$$ LANGUAGE sql IMMUTABLE;