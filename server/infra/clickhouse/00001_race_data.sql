-- +goose Up
CREATE TABLE telemetry.car_motion
(
    session_id       String,
    frame_id         UInt32,
    overall_frame_id UInt32,
    session_time     Float32,
    car_index        UInt8,
    recorded_at      DateTime64(3) DEFAULT now64(3),
    world_pos_x      Float32 CODEC(Gorilla),
    world_pos_y      Float32 CODEC(Gorilla),
    world_pos_z      Float32 CODEC(Gorilla),
    world_vel_x      Float32 CODEC(Gorilla),
    world_vel_y      Float32 CODEC(Gorilla),
    world_vel_z      Float32 CODEC(Gorilla),
    fwd_dir_x        Int16,
    fwd_dir_y        Int16,
    fwd_dir_z        Int16,
    right_dir_x      Int16,
    right_dir_y      Int16,
    right_dir_z      Int16,
    g_force_lat      Float32 CODEC(Gorilla),
    g_force_lon      Float32 CODEC(Gorilla),
    yaw              Float32 CODEC(Gorilla),
    pitch            Float32 CODEC(Gorilla),
    roll             Float32 CODEC(Gorilla)
)
ENGINE = MergeTree
PARTITION BY toDate(recorded_at)
ORDER BY (session_id, car_index, frame_id)
TTL recorded_at + INTERVAL 90 DAY;

-- +goose Down
DROP TABLE telemetry.car_motion;
