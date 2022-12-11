import re
import heapq
from math import lcm
from typing import List, Tuple, Callable


# class representing monkey
class Monkey:
    def __init__(
        self,
        items: List[int],
        operation: Callable[[int], int],
        divisor: int,
        divisible_throw_to: int,
        indivisible_throw_to: int,
    ):
        # worry levels of items held by monkey
        self.items = items
        # operation to be performed on each item
        self.operation = operation
        # divisor to check divisibility for on each item
        self.divisor = divisor
        # monkeys to throw items to if divisible or not divisible
        self.throw_to = {
            True: divisible_throw_to,
            False: indivisible_throw_to,
        }
        # number of inspected items
        self.inspection_count = 0

    def inspect(self, worry_manager: Callable[[int], int]) -> List[Tuple[int, int]]:
        """
        Inspect current items and return list of monkeys to throw items to.
        """
        throw_tos = []

        for item_worry in self.items:
            # perform operation on item worry level
            item_worry = self.operation(item_worry)
            # perform worry management operation on item worry level
            item_worry = worry_manager(item_worry)
            # add to list of tuples that dictate which monkey to throw item to
            throw_tos.append((self.throw_to[item_worry % self.divisor == 0], item_worry))

        n_items = len(self.items)
        # remove inspected items
        for _ in range(n_items):
            self.items.pop(0)

        # increment inspection count
        self.inspection_count += n_items

        return throw_tos


# class for managing list of monkeys
class MonkeysManager:
    def __init__(self):
        # list of monkeys
        self.monkeys = []
        # worry management function
        self.worry_manager = lambda x: x

    def add_monkey(self, monkey: Monkey):
        """
        Add monkey to list.
        """
        self.monkeys.append(monkey)

    def set_worry_manager(self, worry_manager: Callable[[int], int]):
        """
        Set worry management function.
        """
        self.worry_manager = worry_manager

    def run_inspection_round(self):
        """
        Have each monkey inspect items sequentially.
        """
        for monkey in self.monkeys:
            throw_tos = monkey.inspect(worry_manager=self.worry_manager)
            # throw items to other monkeys
            for i, item_worry in throw_tos:
                self.monkeys[i].items.append(item_worry)

    def monkey_business(self, k: int) -> int:
        """
        Compute monkey business using product of top-k number of items inspected.
        """
        # top-k number of items inspected
        top_k_inspection_count = []
        for monkey in self.monkeys:
            # push inspection count to heap while maintaining top-k counts
            if len(top_k_inspection_count) < k:
                heapq.heappush(top_k_inspection_count, monkey.inspection_count)
            else:
                heapq.heappushpop(top_k_inspection_count, monkey.inspection_count)

        # compute product of top-k number of items inspected
        inspections_product = 1
        for i in top_k_inspection_count:
            inspections_product *= i

        return inspections_product


if __name__ == '__main__':
    # file with monkeys input data
    file_path = '../monkeys.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # regex pattern for parsing starting items
    regex_starting_items = re.compile(r'^\s*Starting items:\s*(.+)$')
    # regex pattern for parsing operation
    regex_operation = re.compile(r'^\s*Operation:\s*new\s*=\s*(.+)$')
    # regex pattern for parsing divisor
    regex_divisor = re.compile(r'^\s*Test:\s*divisible\s*by\s*(\d+)$')
    # regex pattern for parsing monkey to throw to if divisible
    regex_divisible_throw_to = re.compile(r'^\s*If\s*true:\s*throw\s*to\s*monkey\s*(\d+)$')
    # regex pattern for parsing monkey to throw to if indivisible
    regex_indivisible_throw_to = re.compile(r'^\s*If\s*false:\s*throw\s*to\s*monkey\s*(\d+)$')

    monkeys_manager = MonkeysManager()
    # list of divisors
    divisors = []

    for monkey_input in file_contents.split('\n\n'):
        lines = monkey_input.splitlines()

        # parse monkey infos
        result = regex_starting_items.search(lines[1])
        starting_items = [int(i.strip()) for i in result.group(1).split(',')]

        result = regex_operation.search(lines[2])
        operation = eval(f'lambda old: {result.group(1)}')

        result = regex_divisor.search(lines[3])
        divisor = int(result.group(1))
        divisors.append(divisor)

        result = regex_divisible_throw_to.search(lines[4])
        divisible_throw_to = int(result.group(1))

        result = regex_indivisible_throw_to.search(lines[5])
        indivisible_throw_to = int(result.group(1))

        # create and add monkey to monkeys manager
        monkey = Monkey(
            items=starting_items,
            operation=operation,
            divisor=divisor,
            divisible_throw_to=divisible_throw_to,
            indivisible_throw_to=indivisible_throw_to,
        )
        monkeys_manager.add_monkey(monkey)

    # compute worry management function
    # worry_manager = lambda x: x // 3
    divisors_lcm = lcm(*divisors)
    worry_manager = lambda x: x % divisors_lcm
    monkeys_manager.set_worry_manager(worry_manager)

    # run inspection rounds
    n_rounds = 10000
    for r in range(n_rounds):
        monkeys_manager.run_inspection_round()

    # compute and print monkey business using top-2 inspection counts
    print(monkeys_manager.monkey_business(2))
