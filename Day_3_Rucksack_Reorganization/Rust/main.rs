use std::fs;
use std::collections::HashSet;

// rucksack contents struct
struct Rucksack {
    // hash sets for rucksack compartment contents
    compartment1: HashSet<char>,
    compartment2: HashSet<char>,
}

impl Rucksack {
    // get duplicate items set between two compartments
    pub fn duplicate_items(&self) -> HashSet<char> {
        let duplicate_items_set: HashSet<char> = self.compartment1.intersection(&self.compartment2).cloned().collect();

        return duplicate_items_set;
    }

    // get set of all items in two compartments combined
    pub fn all_items(&self) -> HashSet<char> {
        let all_items_set: HashSet<char> = self.compartment1.union(&self.compartment2).cloned().collect();

        return all_items_set;
    }
}

// get priority of item based on ASCII value
fn get_item_priority(item: char) -> u8 {
    let item_ascii: u8 = item as u8;
    let lower_a_ascii: u8 = 'a' as u8;
    let upper_a_ascii: u8 = 'A' as u8;
    let lower_z_ascii: u8 = 'z' as u8;

    let item_priority: u8;
    // when item is lowercase, priority is given by how many steps ahead of 'a' the item is
    if item_ascii >= lower_a_ascii {
        item_priority = item_ascii - lower_a_ascii + 1;
    } else {
        // priority is given by how many steps ahead of 'A' the item is
		// plus the number of letters in the alphabet (i.e. 'z' - 'a' + 1)
        item_priority = item_ascii - upper_a_ascii + 1 + lower_z_ascii - lower_a_ascii + 1;
    }

    return item_priority;
}

fn main() {
    // file with input rucksack contents
    let file_path: &str = "../rucksacks.txt";
    // read file contents to string
    let file_contents: String = fs::read_to_string(file_path).expect("failed to read file");
    let mut rucksacks : Vec<Rucksack> = Vec::new();

    for line in file_contents.split("\n") {
        // get half length of line
        let half_len = line.len() / 2;
        // insert first half of line into compartment 1
        let mut compartment1: HashSet<char> = HashSet::new();
        for c in (&line[..half_len]).chars() {
            compartment1.insert(c);
        }
        // insert second half of line into compartment 2
        let mut compartment2: HashSet<char> = HashSet::new();
        for c in (&line[half_len..]).chars() {
            compartment2.insert(c);
        }

        // create rucksack with two compartments
        let rucksack: Rucksack = Rucksack {
            compartment1: compartment1,
            compartment2: compartment2,
        };

        // store ruckstack in vector
        rucksacks.push(rucksack);
    }

    let mut sum_of_priorities: u32 = 0;
    // iterate over rucksacks
    for rucksack in rucksacks.iter() {
        // get set of duplicate items between rucksack compartments
        let duplicate_items_set: HashSet<char> = rucksack.duplicate_items();
        // obtain first item in set as optional
        let first_item_optional: Option<&char> = duplicate_items_set.iter().next();
        match first_item_optional {
            // update sum of item priorities if first item is not None
            Some(item) => sum_of_priorities += get_item_priority(*item) as u32,
            // panic if first item is None
            None => panic!("failed to get common compartment items"),
        }
    }

    println!("{sum_of_priorities}");

    sum_of_priorities = 0;
    // iterate over rucksacks, n rucksacks at a time
    let n: usize = 3;
    for i in (0..rucksacks.len()).step_by(n) {
        // get set of all items in first rucksack's compartments
        // this is the overall set of duplicate items
        let mut duplicate_items_set: HashSet<char> = rucksacks[i].all_items();
        // iterate over next n - 1 rucksacks
        for rucksack in &rucksacks[(i + 1)..(i + n)] {
            // get set of all items in current rucksack's compartments
            let rucksack_all_items: HashSet<char> = rucksack.all_items();
            // retain only what is common between current rucksack and the overall set of duplicate items
            duplicate_items_set = duplicate_items_set.intersection(&rucksack_all_items).cloned().collect();
        }

        // obtain first item in set as optional
        let first_item_optional: Option<&char> = duplicate_items_set.iter().next();
        match first_item_optional {
            // update sum of item priorities if first item is not None
            Some(item) => sum_of_priorities += get_item_priority(*item) as u32,
            // panic if first item is None
            None => panic!("failed to get common compartment items"),
        }
    }

    println!("{sum_of_priorities}");
}