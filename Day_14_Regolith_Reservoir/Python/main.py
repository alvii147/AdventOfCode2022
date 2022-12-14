import numpy as np
from typing import List, Tuple


# exception used specifically to signal a termination in pouring sand
class SandPouringError(BaseException):
    pass


# class representing vertical slice view of cave
class Cave:
    # integers representing air, rock, and sand particles
    AIR = 0
    ROCK = 1
    SAND = 2
    # mode representing part 1 and 2 simulations
    # "abyss" represents part 1, "infinite_floor" represents part 2
    MODES = ['abyss', 'infinite_floor']

    def __init__(self, rocks: List[List[Tuple[int, int]]], sand_j: int, mode: str = 'abyss'):
        max_i = -1
        max_j = -1
        min_j = np.inf

        # iterate over rock paths and obtain cave dimensions
        for path in rocks:
            for rock in path:
                max_i = max(max_i, rock[1])
                max_j = max(max_j, rock[0])
                min_j = min(min_j, rock[0])

        # verify and set mode
        if mode not in self.MODES:
            raise ValueError(f'Invalid mode {mode}')

        self.mode = mode
        # make changes to dimensions for infinite floor mode
        if self.mode == 'infinite_floor':
            # add extra row for floor
            max_i += 1
            # increase horizontal dimension so it's sufficient to contain sand
            max_j += max_i
            min_j -= max_i

        # 2d array representing cave particles
        self.cave = np.zeros((max_i + 1, max_j - min_j + 1), dtype=np.int8)
        # where to pour sand from
        self.sand_j = sand_j - min_j

        # iterate over and draw rock paths in 2d cave array
        for path in rocks:
            for k in range(len(path) - 1):
                self.draw_rock_path(
                    path[k][1],
                    path[k][0] - min_j,
                    path[k + 1][1],
                    path[k + 1][0] - min_j,
                )

    def draw_rock_path(self, i1: int, j1: int, i2: int, j2: int):
        """
        Set path between two coordinates (inclusive) to be filled with rocks.
        """
        # raise error if two coordinates are not on the same row or column
        if i1 != i2 and j1 != j2:
            raise ValueError('Cannot draw rock path diagonally')

        # rearrange ordering so indexing works properly
        if i2 < i1:
            i1, i2 = i2, i1

        # rearrange ordering so indexing works properly
        if j2 < j1:
            j1, j2 = j2, j1

        self.cave[i1 : i2 + 1, j1 : j2 + 1] = self.ROCK

    def pour_sand_abyss(self):
        """
        Get coordinates of next poured sand particle in abyss mode.
        """
        sand_ij = (-1, self.sand_j)

        while True:
            # directly below
            sand_ij_next = (sand_ij[0] + 1, sand_ij[1])
            if sand_ij_next[0] >= self.cave.shape[0]:
                raise SandPouringError

            if self.cave[sand_ij_next] == self.AIR:
                sand_ij = sand_ij_next
                continue

            # bottom left
            sand_ij_next = (sand_ij[0] + 1, sand_ij[1] - 1)
            if sand_ij_next[1] < 0:
                raise SandPouringError

            if self.cave[sand_ij_next] == self.AIR:
                sand_ij = sand_ij_next
                continue

            # bottom right
            sand_ij_next = (sand_ij[0] + 1, sand_ij[1] + 1)
            if sand_ij_next[1] >= self.cave.shape[1]:
                raise SandPouringError

            if self.cave[sand_ij_next] == self.AIR:
                sand_ij = sand_ij_next
                continue

            break

        return sand_ij

    def pour_sand_infinite_floor(self):
        """
        Get coordinates of next poured sand particle in infinite_floor mode.
        """
        if self.cave[0, self.sand_j] == self.SAND:
            raise SandPouringError

        sand_ij = (-1, self.sand_j)

        while True:
            # directly below
            sand_ij_next = (sand_ij[0] + 1, sand_ij[1])
            if sand_ij_next[0] < self.cave.shape[0]:
                if self.cave[sand_ij_next] == self.AIR:
                    sand_ij = sand_ij_next
                    continue

                # bottom left
                sand_ij_next = (sand_ij[0] + 1, sand_ij[1] - 1)
                if sand_ij_next[1] >= 0:
                    if self.cave[sand_ij_next] == self.AIR:
                        sand_ij = sand_ij_next
                        continue

                # bottom right
                sand_ij_next = (sand_ij[0] + 1, sand_ij[1] + 1)
                if sand_ij_next[1] < self.cave.shape[1]:
                    if self.cave[sand_ij_next] == self.AIR:
                        sand_ij = sand_ij_next
                        continue

            break

        if sand_ij[0] == -1:
            raise SandPouringError

        return sand_ij

    def pour_sand_till_filled(self):
        """
        Keep pouring sand particles until filled and return number of particles poured.
        """
        # number of sand particles poured
        count = 0
        while True:
            try:
                # pour sand particle based on mode
                if self.mode == 'abyss':
                    sand_ij = self.pour_sand_abyss()
                elif self.mode == 'infinite_floor':
                    sand_ij = self.pour_sand_infinite_floor()
                else:
                    raise ValueError(f'Invalid mode {self.mode}')

                count += 1
                self.cave[sand_ij] = self.SAND
            except SandPouringError:
                break

        return count

    def __str__(self):
        """
        String representation of cave for debugging.
        """
        chars = {
            self.AIR: '.',
            self.ROCK: '#',
            self.SAND: 'o',
        }

        rows = []
        for i in range(self.cave.shape[0]):
            row = []
            for j in range(self.cave.shape[1]):
                row.append(chars[self.cave[i, j]])

            rows.append(''.join(row))

        return '\n'.join(rows)

    def __repr__(self):
        return self.__str__()


if __name__ == '__main__':
    sand_j = 500
    # file with rocks data
    file_path = '../rocks.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # parse rock paths
    rocks = []
    for line in file_contents.splitlines():
        rocks.append([tuple(int(i) for i in ij.split(',')) for ij in line.split('->')])

    # create cave and pour sand till filled
    cave = Cave(rocks, sand_j, 'infinite_floor')
    print(cave.pour_sand_till_filled())