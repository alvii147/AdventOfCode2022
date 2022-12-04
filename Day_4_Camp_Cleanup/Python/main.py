import re


# parse camp assignmets from input file line
def parse_assignments_line(line):
    # use regex group captures to get camp assignments
    result = re.search(r'^(\d+)-(\d+),(\d+)-(\d+)$', line)
    assignments = [
        [
            int(result.group(1)),
            int(result.group(2)),
        ],
        [
            int(result.group(3)),
            int(result.group(4)),
        ],
    ]

    return assignments


# check if range a contains range b
def range_a_contains_b(a, b):
    return a[0] <= b[0] and a[1] >= b[1]


# check if range a overlaps with range b with a on the left and b on the right
def range_a_left_overlaps_b(a, b):
    return a[0] <= b[1] and a[1] >= b[0]


if __name__ == '__main__':
    # file with camp assignments input
    file_path = '../camp_assignments.txt'
    # open and read file
    with open(file_path) as f:
        lines = f.read().split()

    # parse lines into camp assignments
    camp_assignments = [parse_assignments_line(line) for line in lines]

    # count assignment pairs where one assignment contains another
    contains_count = 0
    for assignments_pair in camp_assignments:
        if range_a_contains_b(*assignments_pair) or range_a_contains_b(*reversed(assignments_pair)):
            contains_count += 1

    print(contains_count)

    # count assignment pairs where one assignment overlaps with another
    overlaps_count = 0
    for assignments_pair in camp_assignments:
        if range_a_left_overlaps_b(*assignments_pair) or range_a_left_overlaps_b(*reversed(assignments_pair)):
            overlaps_count += 1

    print(overlaps_count)
