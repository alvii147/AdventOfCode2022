from __future__ import annotations


class Directory:
    def __init__(self, name: str, parent_directory: Directory = None):
        # name of directory
        self.name = name
        # parent directory object
        self.parent_directory = parent_directory
        # list of immediate subdirectory objects
        self.subdirectories = []
        # dictonary of file names to file sizes
        self.files = {}

    def add_subdirectory(self, subdirectory: Directory):
        """
        Add Directory as a subdirectory.
        """
        if id(self) == id(subdirectory):
            raise ValueError('cannot add self as subdirectory')

        self.subdirectories.append(subdirectory)

    def add_file(self, filename: str, filesize: int):
        """
        Add file with name and size.
        """
        self.files[filename] = filesize

    def get_subdirectory(self, sub_dirname: str) -> Directory:
        """
        Find subdirectory by name.
        """
        # raise exception if at least one matching subdirectory is not found
        try:
            subdirectory = [d for d in self.subdirectories if d.name == sub_dirname][0]
        except IndexError:
            raise ValueError(f'subdirectory with name {sub_dirname} not found')

        return subdirectory

    def get_parent_directory(self) -> Directory:
        """
        Get parent Directory.
        """
        return self.parent_directory

    def compute_size(self: Directory) -> int:
        """
        Recursively compute size of Directory's contents.
        """
        total_size = 0

        # recursively compute sizes of subdirectories
        for subdirectory in self.subdirectories:
            total_size += subdirectory.compute_size()

        # compute sizes of files
        for filesize in self.files.values():
            total_size += filesize

        return total_size


if __name__ == '__main__':
    # file with filesystem data
    file_path = '../filesystem.txt'
    # open and read file
    with open(file_path) as f:
        file_contents = f.read()

    TOTAL_DISK_SPACE = 70000000
    REQUIRED_DISK_SPACE = 30000000
    DIRECTORY_SIZE_THRESHOLD = 100000

    # create root directory
    root = Directory(name='/')
    # list of directories navigated
    navigated_directories = [root]
    # current directory
    pwd = root
    # represents whether or not to expect ls command output
    ls_mode = False

    for line in file_contents.splitlines():
        line_split = line.split()

        # if line is a command (i.e. cd, ls)
        if line_split[0] == '$':
            cmd = line_split[1]

            if cmd == 'cd':
                ls_mode = False
                dirname = line_split[2]

                # navigate to root directory
                if dirname == '/':
                    pwd = root
                # navigate to parent directory
                elif dirname == '..':
                    pwd = pwd.get_parent_directory()
                # navigate to subdirectory
                else:
                    pwd = pwd.get_subdirectory(dirname)
            elif cmd == 'ls':
                # turn on ls mode output
                ls_mode = True
                continue
        # if line is not a command and in ls mode
        elif ls_mode:
            # if directory listed
            if line_split[0] == 'dir':
                dirname = line_split[1]
                # create new directory with current directory as parent
                new_dir = Directory(name=dirname, parent_directory=pwd)
                # add new directory as subdirectory of current directory
                pwd.add_subdirectory(new_dir)
                # add new directory to list of navigated directories
                navigated_directories.append(new_dir)
            # else if file listed
            else:
                filename = line_split[1]
                filesize = int(line_split[0])
                # add file name and size to directory
                pwd.add_file(filename, filesize)
        # raise exception if unknown line encountered
        else:
            raise ValueError(f'encountered non-command input {line} while ls_mode = False')

    # compute sizes of navigated directories
    navigated_directory_sizes = [d.compute_size() for d in navigated_directories]
    # get sum of directory sizes under threshold
    print(sum([size for size in navigated_directory_sizes if size <= DIRECTORY_SIZE_THRESHOLD]))

    # compute size of root directory
    total_size = root.compute_size()
    # compute unused space
    unused_space = TOTAL_DISK_SPACE - total_size
    # compute amount of space needed
    extra_space_needed = REQUIRED_DISK_SPACE - unused_space
    # get size of smallest directory that can be removed to get space needed
    print(min([size for size in navigated_directory_sizes if size >= extra_space_needed]))
