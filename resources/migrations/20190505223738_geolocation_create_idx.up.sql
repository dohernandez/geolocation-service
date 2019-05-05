CREATE UNIQUE INDEX CONCURRENTLY unique_geolocation_idx
ON geolocation (ip_address, country_code, country, city, latitude, longitude, mystery_value);
