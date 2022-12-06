import re
import copy


# parse line in "move x from y to z" format
def parse_move(move):
    # perform regex search
    result = re.search(r'^\s*move\s*(\d+)\s*from\s*(\d*)\s*to\s*(\d*)\s*$', move)
    # capture groups
    # from_idx and to_idx are decremented by 1 from captured groups as indexing begins at 0
    n_crates, from_idx, to_idx = (
        int(result.group(1)),
        int(result.group(2)) - 1,
        int(result.group(3)) - 1,
    )

    return n_crates, from_idx, to_idx


# transfer multiple crates from one stack to another
# retain_order indicates whether or not to move crates in original or reverse order
def transfer_crates(stack_src, stack_dest, n_crates, retain_order=False):
    # get ordering function
    # if retain_order = True, this is a function that returns it's input unchanged
    # if retain_order = False, this is a function that returns a reversed copy of input list
    ordering_func = (lambda x: x) if retain_order else (lambda x: list(reversed(x)))
    # copy over to destination stack
    stack_dest += ordering_func(stack_src[-n_crates:])
    # delete last n crates from source stack
    del stack_src[-n_crates:]


if __name__ == '__main__':
    # file with input crates stacks data
    file_path = '../crates.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # get initial stacks and moves inputs
    stacks_input, moves_input = file_contents.split('\n\n')

    # reverse stack lines so we can parse from bottom up
    stacks_input = list(reversed(stacks_input.split('\n')))
    # count number of total stacks
    n_stacks = len(stacks_input[0].split())

    # initialize n stacks
    stacks = []
    for i in range(n_stacks):
        stacks.append([])

    # remove stack indices
    stacks_input = stacks_input[1:]
    # iterate over stack levels from bottom up
    for level in stacks_input:
        # iterate over every 4th character, starting from index 1
        for j, k in enumerate(range(1, len(level), 4)):
            # append character to appropriate stack if not empty
            if len(level[k].strip()) > 0:
                stacks[j].append(level[k])

    # create copy of stacks
    stacks_copy = copy.deepcopy(stacks)

    # iterate over moves
    for move in moves_input.split('\n'):
        # parse current move
        n_crates, from_idx, to_idx = parse_move(move)
        # transfer appropriate crates
        # set retain_order = False so multiple moved stacks are moved in reverse order
        transfer_crates(stacks[from_idx], stacks[to_idx], n_crates, retain_order=False)

    # join top elements of stacks
    print(''.join([i[-1] for i in stacks]))

    # reset stacks
    stacks = stacks_copy

    # iterate over moves
    for move in moves_input.split('\n'):
        # parse current move
        n_crates, from_idx, to_idx = parse_move(move)
        # transfer appropriate crates
        # set retain_order = True so multiple moved stacks are moved in order
        transfer_crates(stacks[from_idx], stacks[to_idx], n_crates, retain_order=True)

    # join top elements of stacks
    print(''.join([i[-1] for i in stacks]))
