import numpy as np
import heapq
from typing import List


class Grid:
    def __init__(self, input_lines: List[str]):
        self.l = len(input_lines)
        self.w = len(input_lines[0])
        self.elevation = np.zeros((self.l, self.w), dtype=np.int8)
        self.start = (0, 0)
        self.end = (0, 0)
        self.lowest = set()

        for i in range(self.l):
            for j in range(self.w):
                v, is_start, is_end = self.get_elevation(input_lines[i][j])

                if v == 0:
                    self.lowest.add((i, j))

                self.elevation[i, j] = v
                if is_start:
                    self.start = (i, j)
                elif is_end:
                    self.end = (i, j)

    def get_elevation(self, char: str):
        is_start = False
        is_end = False

        elevation_char = char
        if elevation_char == 'S':
            elevation_char = 'a'
            is_start = True
        elif elevation_char == 'E':
            elevation_char = 'z'
            is_end = True

        elevation = ord(elevation_char) - ord('a')

        return elevation, is_start, is_end

    def get_neighbours(self, i, j, reverse=False):
        candidates = [
            (i - 1, j),
            (i + 1, j),
            (i, j - 1),
            (i, j + 1),
        ]
        neighbours = set()

        for _i, _j in candidates:
            if _i < 0 or _i >= self.l or _j < 0 or _j >= self.w:
                continue

            if (_i, _j) in self.visited:
                continue

            if not reverse and self.elevation[_i, _j] - self.elevation[i, j] > 1:
                continue

            if reverse and self.elevation[i, j] - self.elevation[_i, _j] > 1:
                continue

            neighbours.add((_i, _j))

        return neighbours

    def single_source_shortest_path(self):
        self.visited = {self.start}
        queue = [self.start + (0,)]
        while len(queue) > 0:
            i, j, d = queue.pop(0)
            if (i, j) == self.end:
                break

            for _i, _j in self.get_neighbours(i, j):
                self.visited.add((_i, _j))
                queue.append((_i, _j, d + 1))

        return d

    def multiple_sources_shortest_path(self):
        self.visited = {self.end}
        distances = {}
        queue = [(0,) + self.end]
        while len(queue) > 0:
            d, i, j = heapq.heappop(queue)
            if d < distances.get((i, j), float('inf')):
                distances[(i, j)] = d

            for _i, _j in self.get_neighbours(i, j, reverse=True):
                self.visited.add((_i, _j))
                heapq.heappush(queue, (d + 1, _i, _j))

        return distances


if __name__ == '__main__':
    # file with elevation data
    file_path = '../elevation.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    input_lines = file_contents.splitlines()
    grid = Grid(input_lines)
    print(grid.single_source_shortest_path())
    distances = grid.multiple_sources_shortest_path()
    print(min([distances.get((i, j), float('inf')) for (i, j) in grid.lowest]))
