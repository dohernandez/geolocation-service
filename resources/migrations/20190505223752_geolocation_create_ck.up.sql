ALTER TABLE geolocation
ADD CONSTRAINT unique_geolocation_ck
UNIQUE USING INDEX unique_geolocation_idx;
