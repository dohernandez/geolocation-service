CREATE TABLE geolocation (
    id UUID PRIMARY KEY NOT NULL,
    ip_address varchar(15) NOT NULL,
    country_code char(2) NOT NULL,
    country varchar(120) NOT NULL,
    city varchar(120) NOT NULL,
    latitude text NOT NULL,
    longitude text NOT NULL,
    mystery_value bigint NOT NULL,

    created_at timestamp default current_timestamp
);
