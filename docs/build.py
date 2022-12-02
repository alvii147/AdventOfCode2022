import os


SHIELDS_IO_BADGE_URL = 'https://img.shields.io/badge'
LANGUAGE_LOGOS = {
    'Go': 'go-00ADD8?style=for-the-badge&logo=go&logoColor=FFFFFF',
    'Python': 'python-3670A0?style=for-the-badge&logo=python&logoColor=FFDD54',
}


if __name__ == '__main__':
    with open('../README.md', 'w') as readme_file:
        readme_file.write('<p align="center">\n')
        readme_file.write('<img alt="Advent of Code 2022 Logo" src="docs/img/logo.png" width=600 />\n')
        readme_file.write('</p>\n\n')
        readme_file.write('# Advent of Code 2022\n\n')
        readme_file.write(
            'Advent of Code is an Advent calendar of small programming puzzles '
            'for a variety of skill sets and skill levels that can be solved in any programming language you like. '
            'This repository contains solutions to the 2022 Advent of Code calendar.\n\n'
        )

        dirnames = sorted([dirname for dirname in os.listdir('../') if 'day' in dirname.lower()])
        progress_bar = 'https://progress-bar.dev/' + str(round((len(dirnames) / 25) * 100))
        readme_file.write(f'Completed **{len(dirnames)}** out of **25** advent day puzzles.\n\n')
        readme_file.write(f'![Progress Bar]({progress_bar})\n\n')
        readme_file.write('Day | Problem | Languages\n')
        readme_file.write('--- | --- | ---\n')

        for dirname in dirnames:
            day_num = int(dirname.lower().replace('day', ''))
            puzzle_link = f'[See Puzzle](https://adventofcode.com/2022/day/{day_num})'
            languages = [lang for lang in os.listdir(f'../{dirname}') if lang in LANGUAGE_LOGOS]
            language_badges = ' '.join([f'![]({SHIELDS_IO_BADGE_URL}/{LANGUAGE_LOGOS[lang]})' for lang in languages])

            readme_file.write(f'{day_num} | {puzzle_link} | {language_badges}\n')
