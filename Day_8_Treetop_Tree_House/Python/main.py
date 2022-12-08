import numpy as np
import numpy.typing as npt
from typing import Any


def is_cumulative_max(A: npt.NDArray[Any]) -> npt.NDArray[bool]:
    """
    Compute whether or not each array element is the cumulative maximum.
    """
    l = A.shape[0]

    return A[1 : l - 1] > np.maximum.accumulate(A[: l - 2])


def n_first_consecutive(A: npt.NDArray[Any]) -> int:
    """
    Compute number of first consecutive occurrences of True.
    """
    where = np.where(np.invert(A))[0]

    # if all elements are true, return last index
    if where.shape[0] == 0:
        return A.shape[0] - 1

    return where[0]


if __name__ == '__main__':
    # file with trees data
    file_path = '../trees.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # create 2d array of tree heights
    trees = np.array(
        [list(line) for line in file_contents.split() if len(line) > 0],
        dtype=np.int8,
    )
    # get tree grid dimensions
    l, w = trees.shape
    # set of coordinates of non-bordering trees visible from at least one side
    non_bordering_visible = set()

    # iterate over rows
    for i in range(1, l - 1):
        # get trees visible from the left
        visible_from_left = is_cumulative_max(trees[i])
        # get trees visible from the right
        visible_from_right = is_cumulative_max(trees[i, ::-1])

        # perform OR of tree visibilities and update coordinates
        for j in np.where(visible_from_left | visible_from_right[::-1])[0] + 1:
            non_bordering_visible.add((i, j))

    # iterate over columns
    for j in range(1, w - 1):
        # get trees visible from the top
        visible_from_top = is_cumulative_max(trees[:, j])
        # get trees visible from the bottom
        visible_from_bottom = is_cumulative_max(trees[::-1, j])

        # perform OR of tree visibilities and update coordinates
        for i in np.where(visible_from_top | visible_from_bottom[::-1])[0] + 1:
            non_bordering_visible.add((i, j))

    # compute number of trees along the outer border
    # these trees are always visible
    n_bordering_trees = 2 * (l + w - 2)
    # number of all visible trees
    n_visible_trees = len(non_bordering_visible) + n_bordering_trees
    print(n_visible_trees)

    max_scenic_score = 0
    # iterate of rows
    for i in range(l):
        # iterate of columns
        for j in range(w):
            # compute scenic score
            scenic_score = (
                (n_first_consecutive(trees[i, j - 1 :: -1] < trees[i, j]) + 1)
                * (n_first_consecutive(trees[i, j + 1 :] < trees[i, j]) + 1)
                * (n_first_consecutive(trees[i - 1 :: -1, j] < trees[i, j]) + 1)
                * (n_first_consecutive(trees[i + 1 :, j] < trees[i, j]) + 1)
            )

            # update maximum scenic score if current is higher
            if scenic_score > max_scenic_score:
                max_scenic_score = scenic_score

    print(max_scenic_score)
