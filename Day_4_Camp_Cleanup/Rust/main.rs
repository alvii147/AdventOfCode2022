#![allow(non_snake_case)]

#[macro_use]
extern crate lazy_static;

use std::fs;
use regex::Regex;

// check if one range contains another
fn check_range_containment(a1: u32, a2: u32, b1: u32, b2: u32) -> bool {
    return (a1 <= b1 && a2 >= b2) || (b1 <= a1 && b2 >= a2);
}

// check if one range overlaps another
fn check_range_overlap(a1: u32, a2: u32, b1: u32, b2: u32) -> bool {
    return (a1 <= b2 && a2 >= b1) || (b1 <= a2 && b2 >= a1);
}

// parse camp assignmets from input file line
fn parse_assignments_line(line: &str) -> (u32, u32, u32, u32) {
    // compile regex at runtime
    lazy_static! {
        static ref RE: Regex = Regex::new(r"(\d+)-(\d+),(\d+)-(\d+)").unwrap();
    }

    // get captured groups of first match
    let cap: regex::Captures = RE.captures_iter(line).next().unwrap();
    // parse matches into integers
    let (a1, a2, b1, b2): (u32, u32, u32, u32) = (
        cap[1].to_string().parse::<u32>().unwrap(),
        cap[2].to_string().parse::<u32>().unwrap(),
        cap[3].to_string().parse::<u32>().unwrap(),
        cap[4].to_string().parse::<u32>().unwrap(),
    );

    return (a1, a2, b1, b2);
}

fn main() {
    // file with input
    let file_path: &str = "../camp_assignments.txt";
    // read file contents to string
    let file_contents: String = fs::read_to_string(file_path).expect("failed to read file");

    let mut contains_count = 0;
    let mut overlaps_count = 0;

    for line in file_contents.split("\n") {
        // parse lines to get camp assignment ranges
        let (a1, a2, b1, b2): (u32, u32, u32, u32) = parse_assignments_line(line);
        // check for range containment and increment
        contains_count += check_range_containment(a1, a2, b1, b2) as u32;
        // check for range overlap and increment
        overlaps_count += check_range_overlap(a1, a2, b1, b2) as u32;
    }

    println!("{contains_count}");
    println!("{overlaps_count}");
}