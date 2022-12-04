use std::fs;
use std::str::SplitWhitespace;

// scores by outcome
const SCORE_LOSS: u8 = 0;
const SCORE_DRAW: u8 = 3;
const SCORE_WIN: u8 = 6;
// mod shift required to obtain outcome
const SHIFT_LOSS: u8 = 2;
const SHIFT_DRAW: u8 = 0;
const SHIFT_WIN: u8 = 1;

// convert letter to play value
// valid letters are A, B, C, X, Y, & Z
fn letter_to_play_value(letter: char) -> u8 {
    let letter_ascii: u8 = letter as u8;
    let upper_a_ascii: u8 = 'A' as u8;
    let upper_c_ascii: u8 = 'C' as u8;
    let upper_x_ascii: u8 = 'X' as u8;
    let upper_z_ascii: u8 = 'Z' as u8;

    // check that letter is one of A, B, C, X, Y, & Z
    if letter_ascii < upper_a_ascii || (letter_ascii > upper_c_ascii && letter_ascii < upper_x_ascii) || letter_ascii > upper_z_ascii {
        panic!("failed to convert letter {} to play value", letter);
    }

    // compute ascii offset for computing play value
    let offset_ascii: u8 = if letter_ascii < upper_x_ascii { upper_a_ascii } else { upper_x_ascii };
    // convert to play value
    let play_value: u8 = letter_ascii - offset_ascii + 1;

    return play_value;
}

// get rock paper scissors score given opponent and user's plays
// 1 => rock, 2 => paper, 3 => scissors
// score is a sum of user's play and the outcome score
fn rock_paper_scissors_score(opponent_plays: u8, i_play: u8) -> u8 {
    let mut score: u8 = i_play;
    // difference between user's and opponent's play
    let plays_diff: i8 = (i_play as i8) - (opponent_plays as i8);

    // if plays are same, it's a draw
    if plays_diff == 0 {
        score += SCORE_DRAW;
    // if user is 1 ahead or 2 behind, user wins
    } else if plays_diff == 1 || plays_diff == -2 {
        score += SCORE_WIN;
    // otherwise (if opponent is 1 ahead or 2 behind) opponent wins
    } else {
        score += SCORE_LOSS;
    }

    return score;
}

// what to play as user given opponent's play and wanted outcome
fn what_to_play(opponent_plays: u8, outcome: u8) -> u8 {
    // number of steps to shift in order to get user's play
    let shift: u8;
    match outcome {
        1 => shift = SHIFT_LOSS,
        2 => shift = SHIFT_DRAW,
        _ => shift = SHIFT_WIN,
    };

    // compute user's play
    let i_play: u8 = ((opponent_plays + shift - 1) % 3) + 1;

    return i_play;
}

fn main() {
    // file with input rock paper scissors data
    let file_path: &str = "../rockpaperscissors.txt";
    // read file contents to string
    let file_contents: String = fs::read_to_string(file_path).expect("failed to read file");
    // total score on the assumption that the second column represents user's play
    let mut total_score_part_1: u32 = 0;
    // total score on the assumption that the second column represents outcome
    let mut total_score_part_2: u32 = 0;

    for line in file_contents.split("\n") {
        let mut line_split: SplitWhitespace = line.trim().split_whitespace();
        // get opponent's play value
        let opponent_plays: u8 = letter_to_play_value(line_split.next().unwrap().chars().nth(0).unwrap());
        // get user's play value
        let mut i_play: u8 = letter_to_play_value(line_split.next().unwrap().chars().nth(0).unwrap());
        // compute score for round
        total_score_part_1 += rock_paper_scissors_score(opponent_plays, i_play) as u32;

        // set second column value to outcome
        let outcome: u8 = i_play;
        // figure out what user should play
        i_play = what_to_play(opponent_plays, outcome);
        // compute score for round
        total_score_part_2 += rock_paper_scissors_score(opponent_plays, i_play) as u32;
    }

    println!("{total_score_part_1}");
    println!("{total_score_part_2}");
}