import itertools
from functools import cmp_to_key
from typing import List, Any


def compare_packets(v1: int | List[Any], v2: int | List[Any]):
    """
    Compare two given packets recursively. True means v1 comes before v2, False means v1
    comes after v2. None means v1 and v2 are equivalent.
    """
    # if both integers, perform integer comparison
    if isinstance(v1, int) and isinstance(v2, int):
        if v1 == v2:
            return None

        return v1 < v2

    # convert v1 to list if int
    if isinstance(v1, int):
        v1 = [v1]

    # convert v2 to list if int
    if isinstance(v2, int):
        v2 = [v2]

    # shorter of two list lengths
    shorter_len = len(v2)
    # default value to return
    default_return = False

    # default return value should be None if list lengths are equal
    # this means if list lengths are equal AND every list element are equivalent
    # THEN this function should return None
    if len(v1) == len(v2):
        default_return = None

    # set shorter length and default return value for case where v1 is shorter
    if len(v1) < len(v2):
        shorter_len = len(v1)
        default_return = True

    # recursively compare each list element
    for i in range(shorter_len):
        result = compare_packets(v1[i], v2[i])
        # continue with comparison if equivalent
        if result is None:
            continue

        # otherwise return computed result
        return result

    return default_return

if __name__ == '__main__':
    # file with signals data
    file_path = '../signals.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # list of pairs of packets
    packet_pairs = [[eval(packet) for packet in pair.splitlines()] for pair in file_contents.split('\n\n')]

    indices_sum = 0
    # iterate over packet pairs, compare pair, and update sum of indices
    for i, packet_pair in enumerate(packet_pairs):
        if compare_packets(*packet_pair):
            indices_sum += i + 1

    print(indices_sum)

    # flatten list of packets pairs to get list of packets
    packets = list(itertools.chain(*packet_pairs))
    # add divider packets
    packets.append([[2]])
    packets.append([[6]])
    # sort packets using comparison function
    sorted_packets = sorted(
        packets,
        # convert comparison function to key function
        key=cmp_to_key(lambda a, b: -1 if compare_packets(a, b) else 1),
    )

    indices_prod = 1
    # iterate over packets, find dividers, and update product of indices
    for i, packet in enumerate(sorted_packets):
        if packet == [[2]] or packet == [[6]]:
            indices_prod *= (i + 1)

    print(indices_prod)
