DROP TRIGGER IF EXISTS trg_playlists_updated_at ON playlists;
DROP TRIGGER IF EXISTS trg_streams_updated_at   ON streams;
DROP TRIGGER IF EXISTS trg_users_updated_at     ON users;
DROP FUNCTION IF EXISTS touch_updated_at();

DROP TABLE IF EXISTS playlist_tracks;
DROP TABLE IF EXISTS playlists;
DROP TABLE IF EXISTS tracks;
DROP TABLE IF EXISTS streams;
DROP TABLE IF EXISTS users;
