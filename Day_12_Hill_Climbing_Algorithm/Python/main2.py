def get_height(c):
    _c = c
    if c == 'S':
        _c = 'a'
    elif c == 'E':
        _c = 'z'

    return ord(_c) - ord('a')

def get_next_pos(i, j, grid, visited):
    curr_height = get_height(grid[i][j])
    next_pos_candidates = [
        (i - 1, j),
        (i + 1, j),
        (i, j - 1),
        (i, j + 1),
    ]
    next_pos = []

    for _i, _j in next_pos_candidates:
        if _i < 0 or _j < 0:
            continue
        try:
            grid[_i][_j]
        except:
            continue

        if (_i, _j) in visited:
            continue

        new_height = get_height(grid[_i][_j])

        if new_height - curr_height > 1:
            continue

        next_pos.append((_i, _j))

    return next_pos


if __name__ == '__main__':
    # file with input
    file_path = '../input.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    start_pos = []
    grid = []

    i = 0
    for line in file_contents.splitlines():
        row = []
        j = 0
        for c in line:
            if c == 'S' or c == 'a':
                start_pos.append((i, j))
            row.append(c)
            j+= 1
        i += 1

        grid.append(row)

    min_dist = float('inf')
    for si, sj in start_pos:
        visited = {
            (si, sj): True
        }
        queue = [(si, sj, 0)]

        reached = False
        while len(queue) > 0:
            # print('q', queue)
            i, j, distance = queue.pop(0)
            # print(i, j, distance)
            # print(visited)
            if grid[i][j] == 'E':
                reached = True
                break
            # visited[(i, j)] = True
            next_pos = get_next_pos(i, j, grid, visited)
            for _i, _j in next_pos:
                visited[(_i, _j)] = True
                queue.append((_i, _j, distance + 1))

        if reached:
            min_dist = min(min_dist, distance)

    print(min_dist)
