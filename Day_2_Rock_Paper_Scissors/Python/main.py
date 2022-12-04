# scores based on outcome
OUTCOME_SCORES = {
    1: 0,  # loss for user
    2: 3,  # draw
    3: 6,  # win for user
}

# get rock paper scissors score given user's play and opponents's play
# 1 => rock, 2 => paper, 3 => scissors
# score is a sum of user's play and the outcome score
def rock_paper_scissors_score(opponent_plays, i_play):
    # difference between user's and opponent's play
    diff = i_play - opponent_plays

    score = 0
    # if plays are same, it's a draw
    if diff == 0:
        score += OUTCOME_SCORES[2]
    # if user is 1 ahead or 2 behind, user wins
    elif diff == 1 or diff == -2:
        score += OUTCOME_SCORES[3]
    # otherwise (if opponent is 1 ahead or 2 behind) opponent wins
    else:
        score += OUTCOME_SCORES[1]

    score += i_play

    return score


# what to play as user given opponent's play and wanted outcome
def what_to_play(opponent_plays, outcome):
    # outcome to how many steps to shift in order to get user's play
    outcome_mod_shift_map = {
        # loss for user, shift by 2 steps, i.e if opponent plays rock(1), play scissors(rock + 2)
        1: 2,
        # draw, shift by 0 steps, i.e if opponent plays rock(1), play paper(rock + 0)
        2: 0,
        # win for user, shift by 1 step, i.e if opponent plays rock(1), play paper(rock + 1)
        3: 1,
    }

    # compute user's play
    i_play = ((opponent_plays + outcome_mod_shift_map[outcome] - 1) % 3) + 1

    return i_play


if __name__ == '__main__':
    # file with input rock paper scissors data
    file_path = '../rockpaperscissors.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    ascii_a = 97
    ascii_x = 120
    # parse contents into play values
    strategies = [[ord(i[0].lower()) - ascii_a + 1, ord(i[1].lower()) - ascii_x + 1] for i in [i.split() for i in file_contents.split('\n')]]

    total_score = 0
    for strategy in strategies:
        # compute score for round
        total_score += rock_paper_scissors_score(*strategy)

    print(total_score)

    total_score = 0
    for strategy in strategies:
        # determine what to play as user
        i_play = what_to_play(*strategy)
        # compute score for round
        total_score += rock_paper_scissors_score(strategy[0], i_play)

    print(total_score)
