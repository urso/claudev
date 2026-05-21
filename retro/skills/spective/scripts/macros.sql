-- Retro macros for DuckDB
-- Loaded by init-db.sh after table creation

-- Full-text search
INSTALL fts;
LOAD fts;
PRAGMA create_fts_index('msgs', 'uuid', 'message', overwrite=1);

-- Thread traversal macros
CREATE OR REPLACE MACRO thread_back(start_uuid, max_depth := 10) AS TABLE
  WITH RECURSIVE thread AS (
    SELECT uuid, parentUuid, sessionId, timestamp, message, 0 AS depth
    FROM msgs WHERE uuid = start_uuid
    UNION ALL
    SELECT m.uuid, m.parentUuid, m.sessionId, m.timestamp, m.message, t.depth + 1
    FROM msgs m JOIN thread t ON m.uuid = t.parentUuid
    WHERE t.depth < max_depth
  )
  SELECT * FROM thread ORDER BY depth DESC;

CREATE OR REPLACE MACRO thread_forward(start_uuid, max_depth := 20) AS TABLE
  WITH RECURSIVE forward AS (
    SELECT uuid, parentUuid, sessionId, timestamp, message, 0 AS depth
    FROM msgs WHERE uuid = start_uuid
    UNION ALL
    SELECT m.uuid, m.parentUuid, m.sessionId, m.timestamp, m.message, f.depth + 1
    FROM msgs m JOIN forward f ON m.parentUuid = f.uuid
    WHERE f.depth < max_depth
  )
  SELECT * FROM forward ORDER BY depth;

-- Text extraction (filter tool_use/tool_result)
CREATE OR REPLACE MACRO extract_text(content) AS
  CASE
    WHEN json_type(content) = 'ARRAY'
    THEN (
      SELECT string_agg(el.value->>'text', E'\n')
      FROM json_each(content) el
      WHERE el.value->>'type' = 'text'
    )
    ELSE content::TEXT
  END;

-- FTS search with context (returns JSON candidates)
-- Search user messages
CREATE OR REPLACE MACRO search_user(terms, max_results := 50) AS TABLE
  WITH matches AS (
    SELECT uuid, fts_main_msgs.match_bm25(uuid, terms) as score
    FROM msgs
    WHERE score IS NOT NULL AND message.role = 'user'
    ORDER BY score DESC
    LIMIT max_results
  )
  SELECT to_json({
    uuid: m.uuid,
    source: 'user',
    score: matches.score,
    timestamp: m.timestamp,
    project: regexp_extract(m.cwd, '.*/github\\.com/[^/]+/([^/]+)', 1),
    msg: extract_text(m.message.content),
    before: extract_text(p.message.content),
    after: extract_text(c.message.content)
  }) as candidate
  FROM matches
  JOIN msgs m ON matches.uuid = m.uuid
  LEFT JOIN msgs p ON m.parentUuid = p.uuid
  LEFT JOIN msgs c ON c.parentUuid = m.uuid
  WHERE extract_text(m.message.content) NOT LIKE '%/retro%'
    AND extract_text(m.message.content) NOT LIKE '%/init%';

-- Search assistant messages (higher signal for decisions)
CREATE OR REPLACE MACRO search_assistant(terms, max_results := 50) AS TABLE
  WITH matches AS (
    SELECT uuid, fts_main_msgs.match_bm25(uuid, terms) as score
    FROM msgs
    WHERE score IS NOT NULL AND message.role = 'assistant'
    ORDER BY score DESC
    LIMIT max_results
  ),
  -- Filter: require user response within 3 turns (not monologue)
  with_response AS (
    SELECT matches.*, m.uuid as m_uuid
    FROM matches
    JOIN msgs m ON matches.uuid = m.uuid
    WHERE EXISTS (
      SELECT 1 FROM msgs child
      WHERE child.message.role = 'user'
        AND child.sessionId = m.sessionId
        AND child.timestamp > m.timestamp
        AND child.timestamp < m.timestamp + INTERVAL '10 minutes'
    )
  )
  SELECT to_json({
    uuid: m.uuid,
    source: 'assistant',
    score: wr.score,
    timestamp: m.timestamp,
    project: regexp_extract(m.cwd, '.*/github\\.com/[^/]+/([^/]+)', 1),
    msg: extract_text(m.message.content),
    before: extract_text(p.message.content),
    after: extract_text(c.message.content)
  }) as candidate
  FROM with_response wr
  JOIN msgs m ON wr.m_uuid = m.uuid
  LEFT JOIN msgs p ON m.parentUuid = p.uuid
  LEFT JOIN msgs c ON c.parentUuid = m.uuid
  WHERE extract_text(m.message.content) NOT LIKE '%task notification%';

-- Combined search (legacy compatibility)
CREATE OR REPLACE MACRO search_messages(terms, max_results := 50) AS TABLE
  SELECT * FROM search_user(terms, max_results / 2)
  UNION ALL
  SELECT * FROM search_assistant(terms, max_results / 2);

-- Coding sessions: sessions with code edits, grouped by project
DROP TABLE IF EXISTS coding_sessions;
CREATE TABLE coding_sessions AS
SELECT
  sessionId as session_id,
  regexp_extract(filename, 'github.com-([^-]+)-([^/]+)', 1) as org,
  regexp_extract(filename, 'github.com-([^-]+)-([^/]+)', 2) as repo,
  COUNT(*) FILTER (WHERE message.content::VARCHAR LIKE '%"name":"Edit"%') as edits,
  COUNT(*) FILTER (WHERE message.content::VARCHAR LIKE '%"name":"Write"%') as writes,
  list_distinct(flatten(list(
    regexp_extract_all(message.content::VARCHAR, '"file_path":"[^"]+\.([a-z]+)"', 1)
  ))) as extensions,
  MIN(timestamp) as started,
  MAX(timestamp) as ended
FROM msgs
WHERE type = 'assistant'
  AND sessionId IS NOT NULL
GROUP BY sessionId, filename
HAVING edits + writes > 0;

-- Edit sequences: find files with 3+ edits and user messages between
-- Note: min_edits parameter doesn't work in HAVING, so we filter after
CREATE OR REPLACE MACRO edit_sequences(min_edits := 3) AS TABLE
  WITH edits AS (
    SELECT
      uuid,
      sessionId,
      timestamp,
      json_extract_string(tool_use.input, '$.file_path') as file_path
    FROM msgs,
      LATERAL (
        SELECT unnest(from_json(message.content::json, '["json"]')) as tool_use
      )
    WHERE type = 'assistant'
      AND json_extract_string(tool_use, '$.type') = 'tool_use'
      AND json_extract_string(tool_use, '$.name') = 'Edit'
  ),
  edit_groups AS (
    SELECT
      file_path,
      sessionId,
      MIN(uuid) as start_uuid,
      MAX(uuid) as end_uuid,
      MIN(timestamp) as start_ts,
      MAX(timestamp) as end_ts,
      COUNT(*) as edit_count
    FROM edits
    GROUP BY file_path, sessionId
  ),
  with_user_messages AS (
    SELECT
      eg.*,
      COUNT(m.uuid) > 0 as has_user_between
    FROM edit_groups eg
    LEFT JOIN msgs m
      ON m.sessionId = eg.sessionId
      AND m.type = 'user'
      AND m.timestamp > eg.start_ts
      AND m.timestamp < eg.end_ts
    GROUP BY eg.file_path, eg.sessionId, eg.start_uuid, eg.end_uuid,
             eg.start_ts, eg.end_ts, eg.edit_count
  )
  SELECT
    file_path,
    sessionId as session_id,
    start_uuid,
    end_uuid,
    edit_count
  FROM with_user_messages
  WHERE has_user_between = true
    AND edit_count >= min_edits;
