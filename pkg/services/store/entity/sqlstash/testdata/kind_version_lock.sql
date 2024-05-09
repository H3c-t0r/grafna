SELECT "resource_version"
  FROM "kind_version"
  WHERE 1 = 1
    AND "group"         = ?
    AND "group_version" = ?
    AND "resource"      = ?
  FOR UPDATE
;
