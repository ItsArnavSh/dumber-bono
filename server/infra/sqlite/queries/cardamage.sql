-- name: InsertCarDamage :one
INSERT INTO car_damage (
    session_id, car_index, frame_identifier, lap_distance,
    tyres_wear_rl, tyres_wear_rr, tyres_wear_fl, tyres_wear_fr,
    tyres_damage_rl, tyres_damage_rr, tyres_damage_fl, tyres_damage_fr,
    brakes_damage_rl, brakes_damage_rr, brakes_damage_fl, brakes_damage_fr,
    tyre_blisters_rl, tyre_blisters_rr, tyre_blisters_fl, tyre_blisters_fr,
    front_left_wing_damage, front_right_wing_damage, rear_wing_damage,
    floor_damage, diffuser_damage, sidepod_damage,
    drs_fault, ers_fault,
    gear_box_damage, engine_damage,
    engine_mguh_wear, engine_es_wear, engine_ce_wear, engine_ice_wear, engine_mguk_wear, engine_tc_wear,
    engine_blown, engine_seized
) VALUES (
    ?, ?, ?, ?,
    ?, ?, ?, ?,
    ?, ?, ?, ?,
    ?, ?, ?, ?,
    ?, ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?,
    ?, ?,
    ?, ?,
    ?, ?, ?, ?, ?, ?,
    ?, ?
) RETURNING *;

-- name: GetLatestCarDamageByCar :one
SELECT * FROM car_damage
WHERE session_id = ? AND car_index = ?
ORDER BY frame_identifier DESC
LIMIT 1;

-- name: ListCarDamageBySessionAndCar :many
SELECT * FROM car_damage
WHERE session_id = ? AND car_index = ?
ORDER BY frame_identifier ASC;


-- name: GetEngineFailuresBySession :many
SELECT * FROM car_damage
WHERE session_id = ?
    AND (engine_blown = 1 OR engine_seized = 1)
ORDER BY frame_identifier ASC;
