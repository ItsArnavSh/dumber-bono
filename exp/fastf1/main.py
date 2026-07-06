"""
Get circuit info (corner numbers, corner coordinates, angles, marshal sectors)
for a given F1 track using FastF1.

Install first:
    pip install fastf1

Usage:
    python get_circuit_info.py
"""

import fastf1
import pandas as pd

# Optional but recommended: cache downloaded data locally so repeat runs are instant
fastf1.Cache.enable_cache("./f1_cache")


def get_circuit_info(year: int, gp: str, session_type: str = "R"):
    """
    year: e.g. 2024
    gp: Grand Prix name or round number, e.g. "Silverstone", "Monza", "Spa", or 10
    session_type: "R" (Race), "Q" (Qualifying), "FP1", "FP2", "FP3", "S" (Sprint)
    """
    session = fastf1.get_session(year, gp, session_type)
    session.load(telemetry=False, weather=False)  # we only need circuit info, load fast

    circuit_info = session.get_circuit_info()

    print(f"\n=== Circuit Info: {session.event['EventName']} {year} ===\n")

    # Track rotation (degrees) - useful if you're plotting the track and need to orient it
    print(f"Track rotation: {circuit_info.rotation}°\n")

    # Corners: DataFrame with columns like X, Y, Number, Letter, Angle, Distance
    print("--- Corners ---")
    print(circuit_info.corners.to_string(index=False))

    # Marshal lights positions
    print("\n--- Marshal Lights ---")
    print(circuit_info.marshal_lights.to_string(index=False))

    # Marshal sectors
    print("\n--- Marshal Sectors ---")
    print(circuit_info.marshal_sectors.to_string(index=False))

    return circuit_info


def corners_to_json(circuit_info) -> list[dict]:
    """Convert the corners DataFrame into a clean list of dicts, e.g. for
    matching against F1 25's UDP m_lapDistance or building your own track DB."""
    records = []
    for _, row in circuit_info.corners.iterrows():
        records.append(
            {
                "number": int(row["Number"]),
                "letter": row["Letter"]
                if row["Letter"]
                else None,  # e.g. "A" for 3A/3B chicanes
                "x": float(row["X"]),
                "y": float(row["Y"]),
                "angle": float(row["Angle"]),
                "distance_m": float(
                    row["Distance"]
                ),  # distance along the lap - matches m_lapDistance concept
            }
        )
    return records


if __name__ == "__main__":
    info = get_circuit_info(2024, "Austria", "R")

    corners_json = corners_to_json(info)
    print("\n--- As JSON-ready list ---")
    import json

    print(json.dumps(corners_json, indent=2))
