Coordinates = tuple[int, int, int]


def get_adjacent_coordinates(coordinates: Coordinates) -> set[Coordinates]:
    """
    Given coordinates, get set of top, bottom, front, back, left and right coordinates.
    """
    adjacent_coordinates = set([
        (coordinates[0] - 1, coordinates[1], coordinates[2]),
        (coordinates[0] + 1, coordinates[1], coordinates[2]),
        (coordinates[0], coordinates[1] - 1, coordinates[2]),
        (coordinates[0], coordinates[1] + 1, coordinates[2]),
        (coordinates[0], coordinates[1], coordinates[2] - 1),
        (coordinates[0], coordinates[1], coordinates[2] + 1),
    ])

    return adjacent_coordinates


def coordinates_in_range(
    coordinates: Coordinates,
    min_x: int,
    max_x: int,
    min_y: int,
    max_y: int,
    min_z: int,
    max_z: int,
) -> bool:
    """
    Check if coordinates lie in given range.
    """
    in_range = (
        coordinates[0] >= min_x and
        coordinates[0] <= max_x and
        coordinates[1] >= min_y and
        coordinates[1] <= max_y and
        coordinates[2] >= min_z and
        coordinates[2] <= max_z
    )

    return in_range


def exterior_air_coordinates(
    starting_coordinates: Coordinates,
    boulders: set[Coordinates],
    min_x: int,
    max_x: int,
    min_y: int,
    max_y: int,
    min_z: int,
    max_z: int,
) -> set[Coordinates]:
    """
    Get set of coordinates of exterior air (i.e. air not trapped in air pockets).
    """
    # set of exterior air coordinates
    exterior_air = set()
    # set of already visited coordinates
    visited = set([starting_coordinates])
    # coordinates queue
    queue = [starting_coordinates]

    while len(queue) > 0:
        # get next coordinates from queue
        coordinates = queue.pop(0)
        exterior_air.add(coordinates)
        # get surrounding coordinates as candidates for next traversal
        adjacent_coordinates = get_adjacent_coordinates(coordinates)

        for adjacent_coor in adjacent_coordinates:
            # skip if visited
            if adjacent_coor in visited:
                continue

            visited.add(adjacent_coor)

            # skip if blocked by boulders
            if adjacent_coor in boulders:
                continue

            # skip if out of range
            if not coordinates_in_range(
                adjacent_coor,
                min_x,
                max_x,
                min_y,
                max_y,
                min_z,
                max_z,
            ):
                continue

            # add coordinates to queue
            queue.append(adjacent_coor)

    return exterior_air


def get_boulders_range(
    boulders: set[Coordinates],
) -> tuple[int, int, int, int, int, int]:
    """
    Given boulders coordinates, get minimum and maximum coordinates in all dimensions.
    """
    min_x = float('inf')
    max_x = -float('inf')
    min_y = float('inf')
    max_y = -float('inf')
    min_z = float('inf')
    max_z = -float('inf')

    # iterate over boulders and update min/max values
    for boulder in boulders:
        min_x = min(min_x, boulder[0])
        max_x = max(max_x, boulder[0])
        min_y = min(min_y, boulder[1])
        max_y = max(max_y, boulder[1])
        min_z = min(min_z, boulder[2])
        max_z = max(max_z, boulder[2])

    return min_x, max_x, min_y, max_y, min_z, max_z


def naive_surface_area(boulders: set[Coordinates]) -> int:
    """
    Naively compute surface area by counting touching boulder sides.
    """
    surface_area = 0
    for boulder in boulders:
        total_sides = 6
        # get adjacent boulder coordinates
        adjacent_boulders = get_adjacent_coordinates(boulder).intersection(boulders)
        # number of sides covered by boulders
        covered_sides = len(adjacent_boulders)

        # update surface area
        surface_area += total_sides - covered_sides

    return surface_area


def exterior_surface_area(boulders: set[Coordinates]) -> int:
    """
    Compute exterior surface area by counting boulder sides touching exterior air.
    """
    # get iteration range
    min_x, max_x, min_y, max_y, min_z, max_z = get_boulders_range(boulders)
    # get exterior air coordinates
    exterior_air = exterior_air_coordinates(
        (min_x, min_y, min_z),
        boulders,
        min_x,
        max_x,
        min_y,
        max_y,
        min_z,
        max_z,
    )

    surface_area = 0
    for boulder in boulders:
        # get boulder's surrounding coordinates
        adjacent_coordinates = get_adjacent_coordinates(boulder)

        for adjacent_coor in adjacent_coordinates:
            # check if coordinates in range
            in_range = coordinates_in_range(
                adjacent_coor,
                min_x,
                max_x,
                min_y,
                max_y,
                min_z,
                max_z,
            )
            # if coordinates belong to exterior air, then boulder's side is exposed
            # additionally, if coordinates are out of range, boulder's side is exposed
            if adjacent_coor in exterior_air or not in_range:
                surface_area += 1

    return surface_area


if __name__ == '__main__':
    # file with boulder locations
    file_path = '../boulders.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # parse boulder coordinates from file
    boulders = set([
        tuple([int(dim) for dim in line.split(',')])
        for line in file_contents.splitlines()
    ])

    print(naive_surface_area(boulders))
    print(exterior_surface_area(boulders))