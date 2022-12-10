if __name__ == '__main__':
    # file with instructions
    file_path = '../instructions.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    # current cycle
    cycle = 0
    # register
    X = 1
    # signal strength, i.e. the cycle number multiplied by the value of X
    signal_strength = 0
    # CRT writing position
    i = 0
    # length of screen
    screen_length = 40
    # length of sprite, must be odd
    sprite_length = 3
    # initialize row of pixels
    row = [' '] * screen_length

    for line in file_contents.splitlines():
        line_split = line.split()
        instruction = line_split[0]
    
        # single cycle instruction where X is unchanged
        if instruction == 'noop':
            V = 0
            n_cycles = 1
        # double cycle instruction where X is changed
        elif instruction == 'addx':
            V = int(line_split[1])
            n_cycles = 2
        else:
            raise ValueError(f'encountered unknown instruction {instruction}')

        # iterate over cycles
        for _ in range(n_cycles):
            # increment cycle count
            cycle += 1
            # update signal strength 20th cycle or a multiple of 40 cycles after that
            if cycle == 20 or (cycle - 20) % 40 == 0:
                signal_strength += cycle * X

            # check if pixel should be lit or dark
            if i >= X - (sprite_length // 2) and i <= X + (sprite_length // 2):
                row[i] = '#'
            else:
                row[i] = '.'

            # increment writing position
            i += 1
            # if reached end of row
            if i == 40:
                # print row
                print(''.join(row))
                # reset writing position and row
                i = 0
                row = [' '] * 40

        # update register X
        X += V

    print(signal_strength)
