use std::fs;
use std::collections::BinaryHeap;
use std::cmp::Reverse;

fn main() {
    // file with input calories data
    let file_path: &str = "../calories.txt";
    // read file contents to string
    let file_contents: String = fs::read_to_string(file_path).expect("failed to read file");

    // store top-k calories in min heap
    let k: i8 = 3;
    let mut top_k_calories: BinaryHeap<Reverse<u32>> = BinaryHeap::new();
    // initialize min heap with k zeros
    for _ in 0..k {
        top_k_calories.push(Reverse(0));
    }

    // loop over calories for each elf
    for elf_calories in file_contents.split("\n\n") {
        let mut elf_calories_sum: u32 = 0;
        // loop over each item for elf
        for calories in elf_calories.split("\n") {
            elf_calories_sum += calories.parse::<u32>().unwrap();
        }

        // push calories sum onto heap
        top_k_calories.push(Reverse(elf_calories_sum));
        // pop lowest calories off of heap
        top_k_calories.pop();
    }

    let mut top_k_calories_sum: u32 = 0;
    // compute sum over top-k calories
    for _ in 0..k {
        if let Some(Reverse(v)) = top_k_calories.pop() {
            top_k_calories_sum += v;
        }
    }

    println!("{top_k_calories_sum}");
}
