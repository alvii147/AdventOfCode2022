import re


def manhattan_distance(x1: int, y1: int, x2: int, y2: int) -> int:
    """
    Compute Manhattan distance between two coordinates
    """
    return abs(y2 - y1) + abs(x2 - x1)


def tuning_frequency(x: int, y: int, c: int) -> int:
    """
    Compute tuning frequency given coordinates and a constant c. This assumes both x and
    y are within the range [0, c].
    """
    return (c * x) + y


def in_range(v: int, _range: tuple[int, int]) -> bool:
    """
    Check if a value is contained within a given range.
    """
    return v >= _range[0] and v <= _range[1]


def in_ranges(v: int, ranges: list[tuple[int, int]]) -> bool:
    """
    Check if a value is contained within at least one of the given ranges.
    """
    for _range in ranges:
        if in_range(v, _range):
            return True

    return False


def get_exclusion_range(
    y: int,
    sx: int,
    sy: int,
    distance: int,
) -> tuple[int, int] | None:
    """
    Given a sensor location, it's distance to the closest beacon, and row y, get the x
    coordinates of left and right ends of the exclusion row.

    For example, in the diagram below, if y = 3 and distance = 4, the coordinates marked
    "L" and "R" would be returned. "#" represents coordinates where a beacon cannot
    exist.

    y
    0      #
    1     ###
    2    #####
    3   L#####R
    4  ####S####
    5   #######
    6    #####
    7     ###
    8      #
    """
    # return early if y is out of range of beacon exclusion zone
    y_diff = abs(y - sy)
    if y_diff > distance:
        return None

    x_offset = abs(distance - y_diff)
    lx, rx = sx - x_offset, sx + x_offset

    return lx, rx


def beacon_exclusion_row(
    y: int,
    sensors: list[tuple[int, int, int]],
    beacons: set[tuple[int, int]],
) -> int:
    """
    Get number of points along a row where a beacon cannot exist.
    """
    # initialize x coordinate limits
    min_x = float('inf')
    max_x = -float('inf')
    ranges = []

    # iterate over sensors
    for sx, sy, distance in sensors:
        # get exclusion range at row y for
        exclusion_range = get_exclusion_range(y, sx, sy, distance)
        if exclusion_range is None:
            continue

        # add to list of exclusion ranges
        lx, rx = exclusion_range
        ranges.append((lx, rx))

        # update x coordinate limits
        min_x = min(min_x, lx)
        max_x = max(max_x, rx)

    count = 0
    # iterate over row
    # count number of points contained by at least one exclusion range
    for x in range(min_x, max_x + 1):
        # skip check if beacon already exists
        if (x, y) in beacons:
            continue

        if in_ranges(x, ranges):
            count += 1

    return count


def beacon_can_exist(x: int, y: int, sensors: list[tuple[int, int, int]]) -> bool:
    """
    Check if beacon can exist in at given coordinates.
    """
    # iterate over sensors
    for sx, sy, distance in sensors:
        # beacon cannot exist if the distance between it and sensor
        # is less than the distance between the sensor and it's closest beacon
        if manhattan_distance(x, y, sx, sy) <= distance:
            return False

    return True


def find_beacon_location(
    sensors: list[tuple[int, int, int]],
    allowed_range: tuple[int, int],
) -> tuple[int, int] | None:
    """
    Given sensors and a range of coordinates, find coordinates where a beacon can exist.
    This assumes there is exactly one such location in the given range, and will fail to
    provide the right coordinates if this assumption is incorrect.
    """
    # iterate over sensors
    for sx, sy, distance in sensors:
        # iterate over rows within sensor's exclusion zone
        for y in range(sy - distance, sy + distance + 1):
            # skip if y coordinate not in allowed range
            if not in_range(y, allowed_range):
                continue

            # get exclusion range at row y
            lx, rx = get_exclusion_range(y, sx, sy, distance)

            # check if beacon can exist to the left of the exclusion zone
            if in_range(lx - 1, allowed_range) and beacon_can_exist(lx - 1, y, sensors):
                return lx - 1, y

            # check if beacon can exist to the right of the exclusion zone
            if in_range(rx + 1, allowed_range) and beacon_can_exist(rx + 1, y, sensors):
                return rx + 1, y

    return None


if __name__ == '__main__':
    # file with sensors and beacons data
    file_path = '../sensors_and_beacons.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # regex pattern for parsing input
    regex_pattern = re.compile(
        r'^\D*x\s*=\s*(-?\d+)\D*y\s*=\s*(-?\d+)\D*x\s*=\s*(-?\d+)\D*y\s*=\s*(-?\d+)\D*$'
    )

    # list of tuples
    # representing sensor locations and distance to corresponding closest beacon
    sensors = []
    # set of coordinates of beacons
    beacons = set()

    for line in file_contents.splitlines():
        # parse and store sensor and beacon locations
        result = regex_pattern.search(line)
        sx, sy, bx, by = [int(i) for i in result.groups()]
        sensors.append((sx, sy, manhattan_distance(sx, sy, bx, by)))
        beacons.add((bx, by))

    # y = 10
    # allowed_range = (0, 20)
    y = 2000000
    allowed_range = (0, 4000000)
    # constant for finding tuning frequency
    c = 4000000

    beacon_exclusion_row_count = beacon_exclusion_row(y, sensors, beacons)
    print(beacon_exclusion_row_count)

    bx, by = find_beacon_location(sensors, allowed_range)
    freq = tuning_frequency(bx, by, c)
    print(freq)
