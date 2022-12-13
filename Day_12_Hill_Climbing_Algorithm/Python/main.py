import numpy as np
from numpy.typing import NDArray
import heapq
from typing import List, Tuple


# class representing 2d elevation grid
class Grid:
    def __init__(self, input_lines: List[str]):
        # length of grid
        self.l = len(input_lines)
        # width of grid
        self.w = len(input_lines[0])
        # 2d elevation grid
        self.elevation = np.zeros((self.l, self.w), dtype=np.int8)
        # starting point
        self.start = (0, 0)
        # ending point
        self.end = (0, 0)
        # set of points with zero elevation
        self.zero_elevation = set()

        # iterate over grid rows
        for i in range(self.l):
            # iterate over grid columns
            for j in range(self.w):
                # get elevation value from character
                v, is_start, is_end = self.get_elevation(input_lines[i][j])

                # add to zero elevation set if zero
                if v == 0:
                    self.zero_elevation.add((i, j))

                # set elevation value
                self.elevation[i, j] = v

                # set start or end points
                if is_start:
                    self.start = (i, j)
                elif is_end:
                    self.end = (i, j)

    def get_elevation(self, char: str) -> Tuple[int, bool, bool]:
        """
        Get elevation value and whether it's a start or end point, given a character.
        """
        # whether it's a start point
        is_start = False
        # whether it's an end point
        is_end = False

        # a-z elevation character
        elevation_char = char

        # start point
        if elevation_char == 'S':
            elevation_char = 'a'
            is_start = True
        # end point
        elif elevation_char == 'E':
            elevation_char = 'z'
            is_end = True

        # convert character to elevation value
        elevation = ord(elevation_char) - ord('a')

        return elevation, is_start, is_end

    def get_neighbours(self, i: int, j: int, visited: NDArray[np.bool_], reverse: bool = False) -> List[Tuple[int, int]]:
        """
        Get neighbouring indices. This returns only indices of neighbours that can be
        travelled to from the current point. If reverse = True, this returns indices of
        neighbours that can be travelled from to the current point.
        """
        candidates = [
            # top
            (i - 1, j),
            # bottom
            (i + 1, j),
            # left
            (i, j - 1),
            # right
            (i, j + 1),
        ]
        # set of neighbouring indices
        neighbours = set()

        for _i, _j in candidates:
            # not a valid neighbour if indices out of range
            if _i < 0 or _i >= self.l or _j < 0 or _j >= self.w:
                continue

            # not a valid neighbour if already visited
            if visited[_i, _j]:
                continue

            # not a valid neighbour if higher difference in elevation than a single step
            if not reverse and self.elevation[_i, _j] - self.elevation[i, j] > 1:
                continue

            # not a valid neighbour if lower difference in elevation than a single step
            if reverse and self.elevation[i, j] - self.elevation[_i, _j] > 1:
                continue

            # add indices to set of neighbours
            neighbours.add((_i, _j))

        return neighbours

    def multiple_sources_shortest_paths(self, source: Tuple[int, int], reverse: bool) -> NDArray[np.int64]:
        """
        Perform Dikstra's algorithm to obtain shortest paths from given source to
        multiple sources.
        """
        # boolean array representing whether indices have already been visited
        visited = np.zeros((self.l, self.w), dtype=np.bool_)
        visited[source] = True

        # shortest distances from each pair of indices, initialized to inf
        distances = np.full((self.l, self.w), np.iinfo(np.int64).max, dtype=np.int64)

        # priority queue with distance and indices
        # each item in queue is a 3-element tuple
        # first element is distance, so priority queue is always ordered on distance
        queue = [(0,) + source]

        while len(queue) > 0:
            # pop queue item
            d, i, j = heapq.heappop(queue)
            # relaxation
            if d < distances[i, j]:
                distances[i, j] = d

            # push neighbours onto queue
            for _i, _j in self.get_neighbours(i, j, visited, reverse=reverse):
                visited[_i, _j] = True
                heapq.heappush(queue, (d + 1, _i, _j))

        return distances


if __name__ == '__main__':
    # file with elevation data
    file_path = '../elevation.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # create grid by parsing file contents
    grid = Grid(file_contents.splitlines())
    # compute shortest paths to end point
    distances = grid.multiple_sources_shortest_paths(grid.end, True)
    # shortest path from start to end
    print(distances[grid.start])
    # shortest path from any zero elevation point to end
    print(min([distances[i, j] for (i, j) in grid.zero_elevation]))
