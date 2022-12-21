import numpy as np
import numpy.typing as npt


def array_first_true(A: npt.NDArray[np.bool_]) -> int:
    """
    Get index of first true value in array. Returns length of array if none found.
    """
    if not A.any():
        return A.shape[0]

    return np.argmax(A)


# class representing a rock and its structure
class Rock:
    def __init__(
        self,
        dimensions:
        tuple[int, int],
        structure: tuple[tuple[int, int],...],
    ):
        self.height, self.width = dimensions
        # tuple of coordinates with respect to top left
        self.structure = structure

    def __str__(self) -> str:
        """
        String representation of rock for debugging.
        """
        matrix = [['.'] * self.width for _ in range(self.height)]
        for i, j in self.structure:
            matrix[i][j] = '#'

        rows = []
        for row in matrix:
            rows.append(''.join(row))

        return '\n'.join(rows)


# class representing vertical tunnel of rocks
class Tunnel:
    def __init__(
        self,
        jet_patterns: list,
        width: int = 7,
        left_offset: int = 2,
        bottom_offset: int = 3,
    ):
        # list of jet patterns strings
        # can be either '<' or '>'
        self.jet_patterns = jet_patterns
        # current jet pattern index
        self.jet_idx = 0
        # width of tunnel
        self.width = width
        # how far from the left wall each rock stars falling from
        self.left_offset = left_offset
        # how far from the bottom each rock stars falling from
        self.bottom_offset = bottom_offset
        # tunnel grid of rocks
        self.grid = np.zeros((0, self.width), dtype=np.bool_)
        # amount of height that has been truncated to save memory
        self.truncated_height = 0
        # cache for detecting cycles
        self.cycle_cache = {}
        # number of rocks dropped
        self.rock_count = 0

    def __str__(self) -> str:
        """
        String representation of tunnel for debugging.
        """
        return '\n'.join([''.join(row) for row in np.where(self.grid, '#', '.')])

    def get_jet_pattern(self) -> str:
        """
        Get next jet pattern string.
        """
        jet_pattern = self.jet_patterns[self.jet_idx]
        self.jet_idx = (self.jet_idx + 1) % len(self.jet_patterns)

        return jet_pattern

    def detect_collision(self, rock: Rock, top_left: tuple[int, int]):
        """
        Given top left coordinates of rock, detect collision with walls or other rocks.
        This raises CollisionException in the event of a collision.
        """
        # mask array to keep track of rocks
        mask = np.zeros(self.grid.shape, dtype=np.bool_)

        # iterate over rock structure
        for i, j in rock.structure:
            # rock structure coordinates are with respect to top left
            # get global coordinates of structure, with respect to tunnel
            gi, gj = top_left[0] + i, top_left[1] + j

            # detect collision with walls
            if gi < 0 or gi >= self.grid.shape[0] or gj < 0 or gj >= self.width:
                raise CollisionException

            mask[gi, gj] = True

        # detect collision with other rocks
        if (mask & self.grid).any():
            raise CollisionException

    def place_rock(self, rock: Rock, top_left: tuple[int, int]):
        """
        Place rock at given coordinates.
        """
        # iterate over rock structure
        for i, j in rock.structure:
            # rock structure coordinates are with respect to top left
            # get global coordinates of structure, with respect to tunnel
            gi, gj = top_left[0] + i, top_left[1] + j

            # place rock in tunnel grid
            self.grid[gi, gj] = True

    def truncate_grid(self):
        """
        Truncate tunnel from the top and bottom.
        """
        # determine truncation indices
        top = np.argmax(self.grid.any(axis=1))
        bottom = np.amax(np.apply_along_axis(array_first_true, 0, self.grid))
        # keep track of truncated height
        self.truncated_height += self.grid.shape[0] - bottom
        # truncate height
        self.grid = self.grid[top:bottom]

    def drop_rock(self, rock: Rock) -> tuple[bytes, int, int] | None:
        """
        Drop given rock into tunnel and place it where the rock stops moving.
        """
        # cache state of grid, rock used, and the jet pattern index
        cache_item = (self.grid.tobytes(), id(rock), self.jet_idx)
        # return cached item if already cached
        if cache_item in self.cycle_cache:
            return self.cycle_cache[cache_item]

        # cache current state
        self.cycle_cache[cache_item] = self.rock_count
        self.rock_count += 1

        # increase height of tunnel before adding rock
        self.grid = np.vstack(
            (
                np.zeros((rock.height, self.width), dtype=np.bool_),
                np.zeros((self.bottom_offset, self.width), dtype=np.bool_),
                self.grid,
            )
        )

        # loop until collision occurs
        top_left = (0, self.left_offset)
        collision = False
        while not collision:
            jet_pattern = self.get_jet_pattern()
            if jet_pattern == '<':
                next_top_left = (top_left[0], top_left[1] - 1)
            elif jet_pattern == '>':
                next_top_left = (top_left[0], top_left[1] + 1)

            # allow collisions with walls
            # collision with walls means the rock stays in the same place
            try:
                self.detect_collision(rock, next_top_left)
                top_left = next_top_left
            except CollisionException:
                pass

            # collision with other rocks means the rock rests
            next_top_left = (top_left[0] + 1, top_left[1])
            try:
                self.detect_collision(rock, next_top_left)
                top_left = next_top_left
            except CollisionException:
                collision = True

        self.place_rock(rock, top_left)
        self.truncate_grid()

        return None

    def get_height(self) -> int:
        """
        Get current height of tunnel.
        """
        return self.grid.shape[0] + self.truncated_height


# exception to indicate collisions
class CollisionException(BaseException):
    pass


if __name__ == '__main__':
    # file with jet patterns
    file_path = '../jet_patterns.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # get list of jet patterns
    jet_patterns = list(file_contents.strip())
    # create rocks with given structures
    rock_patterns = [
        #   Rock Structure:
        #   ####
        Rock(
            dimensions=(1, 4),
            structure=(
                (0, 0),
                (0, 1),
                (0, 2),
                (0, 3),
            )
        ),
        #   Rock Structure:
        #   .#.
        #   ###
        #   .#.
        Rock(
            dimensions=(3, 3),
            structure=(
                (0, 1),
                (1, 0),
                (1, 1),
                (1, 2),
                (2, 1),
            )
        ),
        #   Rock Structure:
        #   ..#
        #   ..#
        #   ###
        Rock(
            dimensions=(3, 3),
            structure=(
                (0, 2),
                (1, 2),
                (2, 0),
                (2, 1),
                (2, 2),
            )
        ),
        #   Rock Structure:
        #   #
        #   #
        #   #
        #   #
        Rock(
            dimensions=(4, 1),
            structure=(
                (0, 0),
                (1, 0),
                (2, 0),
                (3, 0),
            )
        ),
        #   Rock Structure:
        #   ##
        #   ##
        Rock(
            dimensions=(2, 2),
            structure=(
                (0, 0),
                (0, 1),
                (1, 0),
                (1, 1),
            )
        ),
    ]

    t = Tunnel(jet_patterns)
    n_rock_patterns = len(rock_patterns)
    # number of rocks to drop
    n_rocks = 1000000000000
    # list of heights after every rock is dropped
    heights = []
    for i in range(n_rocks):
        starting_idx = t.drop_rock(rock_patterns[i % n_rock_patterns])
        if starting_idx is not None:
            break

        heights.append(t.get_height())

    cycle_length = i - starting_idx
    starting_height = heights[starting_idx - 1]
    cycle_height = heights[i - 1] - heights[starting_idx - 1]
    n_cycles = (n_rocks - starting_idx) // cycle_length
    remainder_height = heights[n_rocks - (n_cycles * cycle_length) - 1] - starting_height

    print(starting_height + remainder_height + (n_cycles * cycle_height))
