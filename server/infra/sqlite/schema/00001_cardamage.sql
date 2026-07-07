-- +goose Up
CREATE TABLE car_damage (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id               TEXT    NOT NULL,
    car_index                INTEGER NOT NULL,
    frame_identifier          INTEGER NOT NULL,
    lap_distance              REAL,

    tyres_wear_rl             REAL    NOT NULL,
    tyres_wear_rr             REAL    NOT NULL,
    tyres_wear_fl             REAL    NOT NULL,
    tyres_wear_fr             REAL    NOT NULL,

    tyres_damage_rl           INTEGER NOT NULL,
    tyres_damage_rr           INTEGER NOT NULL,
    tyres_damage_fl           INTEGER NOT NULL,
    tyres_damage_fr           INTEGER NOT NULL,

    brakes_damage_rl          INTEGER NOT NULL,
    brakes_damage_rr          INTEGER NOT NULL,
    brakes_damage_fl          INTEGER NOT NULL,
    brakes_damage_fr          INTEGER NOT NULL,

    tyre_blisters_rl          INTEGER NOT NULL,
    tyre_blisters_rr          INTEGER NOT NULL,
    tyre_blisters_fl          INTEGER NOT NULL,
    tyre_blisters_fr          INTEGER NOT NULL,

    front_left_wing_damage    INTEGER NOT NULL,
    front_right_wing_damage   INTEGER NOT NULL,
    rear_wing_damage          INTEGER NOT NULL,
    floor_damage              INTEGER NOT NULL,
    diffuser_damage           INTEGER NOT NULL,
    sidepod_damage            INTEGER NOT NULL,

    drs_fault                 INTEGER NOT NULL, -- 0/1
    ers_fault                 INTEGER NOT NULL, -- 0/1

    gear_box_damage           INTEGER NOT NULL,
    engine_damage             INTEGER NOT NULL,

    engine_mguh_wear          INTEGER NOT NULL,
    engine_es_wear            INTEGER NOT NULL,
    engine_ce_wear            INTEGER NOT NULL,
    engine_ice_wear           INTEGER NOT NULL,
    engine_mguk_wear          INTEGER NOT NULL,
    engine_tc_wear            INTEGER NOT NULL,

    engine_blown              INTEGER NOT NULL, -- 0/1
    engine_seized             INTEGER NOT NULL, -- 0/1

    created_at                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX idx_car_damage_session_car ON car_damage (session_id, car_index);
CREATE INDEX idx_car_damage_session_frame ON car_damage (session_id, frame_identifier);

-- +goose Down
DROP TABLE car_damage;
