import re
import heapq
from itertools import combinations


def valves_floyd_warshall(
    valves_data: dict[str, dict[str, int | frozenset | dict[str, int]]],
):
    """
    Perform Floyd Warshall algorithm to obtain and store distances between all possible
    pairs of valves.
    """
    # iterate over valves
    for valve in valves_data.keys():
        # set distance to self to zero
        valves_data[valve]['distances'] = {
            valve: 0,
        }

        # set distance to neighbours to one
        for neighbouring_valve in valves_data[valve]['neighbours']:
            valves_data[valve]['distances'][neighbouring_valve] = 1

    # compute distances through intermediate valves
    for intermediate_valve in valves_data.keys():
        # iterate over all valve pairs
        for valve1, valve2 in combinations(valves_data.keys(), r=2):
            intermediate_distance = (
                valves_data[valve1]['distances'].get(intermediate_valve, float('inf')) +
                valves_data[intermediate_valve]['distances'].get(valve2, float('inf'))
            )

            current_distance = valves_data[valve1]['distances'].get(valve2, float('inf'))

            if intermediate_distance < current_distance:
                valves_data[valve1]['distances'][valve2] = intermediate_distance
                valves_data[valve2]['distances'][valve1] = intermediate_distance


def max_pressure_released(
    source_valve: str,
    valves_data: dict[str, dict[str, int | frozenset | dict[str, int]]],
    total_minutes: int,
):
    # initiate queue with source valve
    queue = [
        (
            # negative value of total pressure released so far
            # negated because heapq implements min heap
            0,
            # minutes passed
            0,
            # total pressure released so far
            0,
            # current valve
            source_valve,
            # set of unopened valves
            frozenset(valves_data.keys()),
        ),
    ]
    # initalize maximum pressure released
    max_pressure = 0
    # set of valves with zero flow rate
    zero_valves = frozenset([k for k, v in valves_data.items() if v['flow'] < 1])

    while len(queue) > 0:
        # pop from queue
        (
            _,
            minutes_passed,
            pressure,
            current_valve,
            unopened_valves,
        ) = heapq.heappop(queue)
        # number of items pushed to queue in current loop
        push_count = 0

        # iterate over unopened valves with a non-zero flow rate
        for unopened_valve in unopened_valves.difference(zero_valves):
            # minutes needed to travel to and open valve
            minutes_needed = valves_data[current_valve]['distances'][unopened_valve] + 1
            # skip if opening valve will cause time to run out
            if total_minutes - minutes_passed - minutes_needed < 1:
                continue

            # compute amount of pressure released in opening valve
            opening_pressure = (
                # number of minutes valve will flow after opened * flow rate
                (total_minutes - minutes_passed - minutes_needed) *
                valves_data[unopened_valve]['flow']
            )
            # push to queue
            heapq.heappush(
                queue,
                (
                    -(pressure + opening_pressure),
                    minutes_passed + minutes_needed,
                    pressure + opening_pressure,
                    unopened_valve,
                    unopened_valves.difference([unopened_valve]),
                ),
            )
            # increment push count
            push_count += 1

        # update maximum pressure released if no new items push to queue
        if push_count < 1:
            if pressure > max_pressure:
                max_pressure = pressure

    return max_pressure


def max_pressure_released_with_elephant(
    source_valve: str,
    valves_data: dict[str, dict[str, int | frozenset | dict[str, int]]],
    total_minutes: int,
):
    # initiate queue with source valve
    queue = [
        (
            # negative value of total pressure released so far
            # negated because heapq implements min heap
            0,
            # minutes passed for me
            0,
            # minutes passed for elephant
            0,
            # total pressure released so far
            0,
            # current valve for me
            source_valve,
            # current valve for elephant
            source_valve,
            # set of unopened valves
            frozenset(valves_data.keys()),
            # valve opening order path for me
            tuple(),
            # valve opening order path for elephant
            tuple(),
        )]
    # initalize maximum pressure released
    max_pressure = 0
    # set of valves with zero flow rate
    zero_valves = frozenset([k for k, v in valves_data.items() if v['flow'] < 1])
    # set of visited paths
    # visited = set()
    visited = set([(tuple(), tuple())])

    while len(queue) > 0:
        # pop from queue
        (
            _,
            minutes_passed_me,
            minutes_passed_elephant,
            pressure,
            current_valve_me,
            current_valve_elephant,
            unopened_valves,
            path_me,
            path_elephant,
        ) = heapq.heappop(queue)
        # number of items pushed to queue in current loop
        push_count = 0

        # iterate over unopened valves with a non-zero flow rate
        for unopened_valve in unopened_valves.difference(zero_valves):
            # minutes needed to travel to and open valve for me
            minutes_needed = (
                valves_data[current_valve_me]['distances'][unopened_valve] + 1
            )
            # skip if opening valve will cause time to run out
            if total_minutes - minutes_passed_me - minutes_needed < 1:
                continue

            # next path for me
            next_path_me = path_me + (unopened_valve,)
            # skip if path already explored
            if (
                (next_path_me, path_elephant) in visited or
                (path_elephant, next_path_me) in visited
            ):
                continue

            # compute amount of pressure released in opening valve for me
            opening_pressure = (
                # number of minutes valve will flow after opened * flow rate
                (total_minutes - minutes_passed_me - minutes_needed) *
                valves_data[unopened_valve]['flow']
            )
            # add path to already explored set
            visited.add((next_path_me, path_elephant))
            # push to queue
            heapq.heappush(
                queue,
                (
                    -(pressure + opening_pressure),
                    minutes_passed_me + minutes_needed,
                    minutes_passed_elephant,
                    pressure + opening_pressure,
                    unopened_valve,
                    current_valve_elephant,
                    unopened_valves.difference([unopened_valve]),
                    next_path_me,
                    path_elephant,
                ),
            )
            # increment push count
            push_count += 1

        # iterate over unopened valves with a non-zero flow rate
        for unopened_valve in unopened_valves.difference(zero_valves):
            # minutes needed to travel to and open valve for elephant
            minutes_needed = (
                valves_data[current_valve_elephant]['distances'][unopened_valve] + 1
            )
            # skip if opening valve will cause time to run out
            if total_minutes - minutes_passed_elephant - minutes_needed < 1:
                continue

            # next path for elephant
            next_path_elephant = path_elephant + (unopened_valve,)
            # skip if path already explored
            if (
                (path_me, next_path_elephant) in visited or
                (next_path_elephant, path_me) in visited
            ):
                continue

            # compute amount of pressure released in opening valve for elephant
            opening_pressure = (
                # number of minutes valve will flow after opened * flow rate
                (total_minutes - minutes_passed_elephant - minutes_needed) *
                valves_data[unopened_valve]['flow']
            )
            # add path to already explored set
            visited.add((path_me, next_path_elephant))
            # push to queue
            heapq.heappush(
                queue,
                (
                    -(pressure + opening_pressure),
                    minutes_passed_me,
                    minutes_passed_elephant + minutes_needed,
                    pressure + opening_pressure,
                    current_valve_me,
                    unopened_valve,
                    unopened_valves.difference([unopened_valve]),
                    path_me,
                    next_path_elephant,
                ),
            )
            # increment push count
            push_count += 1

        # update maximum pressure released if no new items push to queue
        if push_count < 1:
            if pressure >= max_pressure:
                max_pressure = pressure

    return max_pressure


if __name__ == '__main__':
    # file with valves data
    file_path = '../valves.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # regex pattern for parsing input
    regex_pattern = re.compile(
        r'^Valve\s+(\D{2})\s+has\s+flow\s+rate\s*=\s*(\d+)\s*;\s*tunnels?\s+leads?\s+to\s+valves?\s+([\D,\s]+)$'
    )

    # dictionary for storing valves data
    valves_data = {}

    # iterate over file lines, parse and store contents
    for line in file_contents.splitlines():
        result = regex_pattern.search(line)

        valves_data[result.group(1)] = {
            'flow': int(result.group(2)),
            'neighbours': frozenset([v.strip() for v in result.group(3).split(',')]),
        }

    # compute all valve pair distances
    valves_floyd_warshall(valves_data)

    # print(max_pressure_released('AA', valves_data, 30))
    print(max_pressure_released_with_elephant('AA', valves_data, 26))
