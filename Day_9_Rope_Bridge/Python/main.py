from typing import Tuple


# object representing knots of a rope
class Rope:
    def __init__(self, length: int):
        # length of the rope
        self.length = length
        # initialize coordinates of knots
        self.knots = [(0, 0)] * self.length
        # list of sets of visited coordinates for each knot
        self.visited = []
        for _ in range(self.length):
            self.visited.append(set())

    def sign(self, x: int) -> int:
        """
        Get sign of value.
        """
        if x == 0:
            return 0

        return x // abs(x)

    def follow(self, x2: int, y2: int, x1: int, y1: int) -> Tuple[int, int]:
        """
        Given coordinates of a knot and its following knot, get movement direction for
        following knot.
        """
        delx, dely = x2 - x1, y2 - y1

        # don't move if one step or less away in both dimensions
        if abs(delx) <= 1 and abs(dely) <= 1:
            return 0, 0

        return self.sign(delx), self.sign(dely)

    def move(self, delx: int, dely: int):
        """
        Move first knot by movement direction and update following knots.
        """
        # iterate over knots
        for i in range(self.length):
            # coordinates of current knot
            x1, y1 = self.knots[i]

            # if not head knot
            if i != 0:
                # get movement direction based on previous knot
                delx, dely = self.follow(x2, y2, x1, y1)

            # update knot coordinates
            x1 += delx
            y1 += dely
            self.knots[i] = (x1, y1)

            # update visited coordinates for knot
            self.visited[i].add((x1, y1))

            # set current knot coordinates as previous knot coordinates
            x2, y2 = x1, y1


if __name__ == '__main__':
    # dictionary for converting direction letter into movement direction
    DIRECTIONS_MAP = {
        'U': (0, 1),
        'D': (0, -1),
        'R': (1, 0),
        'L': (-1, 0),
    }
    # length of rope
    length = 10
    # create role of given length
    rope = Rope(length=length)

    # file with rope motions data
    file_path = '../rope_motions.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    for motion in file_contents.splitlines():
        direction, steps = motion.split()
        steps = int(steps)

        # get movement direction
        delx, dely = DIRECTIONS_MAP[direction]

        # move rope step-by-step
        for step in range(steps):
            rope.move(delx, dely)

    # get number of visited coordinates for tail knot
    print(len(rope.visited[length - 1]))
