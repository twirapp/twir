-- Channel API key
ALTER TABLE channels ADD COLUMN api_key TEXT DEFAULT uuidv7();
CREATE UNIQUE INDEX channels_api_key_idx ON channels(api_key) WHERE api_key IS NOT NULL;

-- Hide on pause setting for song requests
ALTER TABLE channels_song_requests_settings ADD COLUMN hide_on_pause BOOL DEFAULT true;
