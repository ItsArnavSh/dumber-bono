CREATE TABLE IF NOT EXISTS telemetry_frames (
    session_id      UINTEGER,
    car_no          UTINYINT,
    frame_time      TIMESTAMP,

    speed           FLOAT,
    throttle        FLOAT,
    steer           FLOAT,
    brake           FLOAT,
    clutch          FLOAT,
    gear            TINYINT,
    engine_rpm      USMALLINT,
    drs             BOOLEAN,

    pos_x           FLOAT,
    pos_y           FLOAT,
    pos_z           FLOAT,

    vel_x           FLOAT,
    vel_y           FLOAT,
    vel_z           FLOAT,

    fwd_x           FLOAT,
    fwd_y           FLOAT,
    fwd_z           FLOAT,

    g_force_lat     FLOAT,
    g_force_lon     FLOAT,

    yaw             FLOAT,
    pitch           FLOAT,
    roll            FLOAT
);

CREATE INDEX IF NOT EXISTS idx_telemetry_session_car_time
    ON telemetry_frames (session_id, car_no, frame_time);
