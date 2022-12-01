import heapq


if __name__ == '__main__':
    # file with input calories data
    file_path = '../calories.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # parse file contents into list of total calories for each elf
    elf_calories = [sum([int(j) for j in i.split()]) for i in file_contents.split('\n\n')]

    # compute and print sum of top 3 calories
    k = 3
    top_k_calories_sum = sum(heapq.nlargest(k, elf_calories))
    print(top_k_calories_sum)
