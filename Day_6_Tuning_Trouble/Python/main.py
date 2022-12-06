from collections import Counter


if __name__ == '__main__':
    # file with input signal
    file_path = '../signal.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # get signal string
    signal = file_contents.strip()
    # number of unique characters to look for
    window_size = 14
    # character counter
    counter = Counter(signal[:window_size])

    # if first window is completely unique, print window size and exit
    if len(counter) == window_size:
        print(window_size)
        exit()

    # iterate over windows
    for i in range(window_size, len(signal) - window_size + 1):
        # increment counter for character to the right of the window
        counter[signal[i]] += 1
        # decrement counter for character to the left of the window
        counter[signal[i - window_size]] -= 1
        if counter[signal[i - window_size]] == 0:
            del counter[signal[i - window_size]]

        # if desired number of unique characters found, print and break
        if len(counter) == window_size:
            print(i + 1)
            break
