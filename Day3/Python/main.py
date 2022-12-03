# get priority of item based on ASCII value
def get_item_priority(item):
    item_ascii = ord(item)
    a_ascii = ord('a')
    A_ascii = ord('A')
    z_ascii = ord('z')

    # when item is lowercase, priority is given by how many steps ahead of 'a' the item is
    item_priority = item_ascii - a_ascii + 1
    # when item is uppercase
    if item_ascii < a_ascii:
        # priority is given by how many steps ahead of 'A' the item is
        item_priority = item_ascii - A_ascii + 1 + z_ascii - a_ascii + 1

    return item_priority


# identify duplicate item between two given strings
def get_duplicate_items(items1, items2):
    # set of elements in items1
    items1_set = set(items1)
    # set of duplicate elements
    duplicate_items = set()

    # iterate over items2 and store any duplicate items
    for item in items2:
        if item in items1_set:
            duplicate_items.add(item)

    # create string of duplicate items
    return ''.join(duplicate_items)

if __name__ == '__main__':
    # file with input rucksack contents
    file_path = '../rucksacks.txt'
    # open and read file
    with open(file_path) as f:
        rucksacks = f.read().split()

    sum_of_priorities = 0
    for rucksack in rucksacks:
        # split rucksack items halfway to get compartment contents
        half_len = len(rucksack) // 2
        compartment1, compartment2 = rucksack[:half_len], rucksack[half_len:]

        # get duplicate items between two compartments
        duplicate_items = get_duplicate_items(compartment1, compartment2)

        # exit early if not exactly one duplicate is found
        if len(duplicate_items) != 1:
            print(f'{len(duplicate_items)} duplicates found in rucksack {rucksack}, expected exactly 1')
            exit()

        # update sum of item priorities
        sum_of_priorities += get_item_priority(duplicate_items)

    print(sum_of_priorities)

    sum_of_priorities = 0
    n = 3
    # iterate over n rucksacks at a time
    for i in range(0, len(rucksacks), n):
        # get n rucksacks
        n_rucksacks = rucksacks[i : i + n]
        # break if less than n rucksacks remaining
        if len(n_rucksacks) != 3:
            break

        duplicate_items = n_rucksacks[0]
        # iterate over rucksack items and gather duplicate items
        for j in range(1, n):
            duplicate_items = get_duplicate_items(duplicate_items, n_rucksacks[j])

        # exit early if not exactly one duplicate is found
        if len(duplicate_items) != 1:
            print(f'{len(duplicate_items)} duplicates found in rucksacks {n_rucksacks}, expected exactly 1')
            exit()

        # update sum of item priorities
        sum_of_priorities += get_item_priority(duplicate_items)

    print(sum_of_priorities)
