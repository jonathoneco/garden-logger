import re
import os

def migrate_files_and_directories(directory):
    """Migrate files and directories from 'n - ' to 'nn. ' pattern."""
    pattern = re.compile(r'^(\d+)\s-\s')

    for root, dirs, files in os.walk(directory, topdown=False):
        # Rename files
        for file in files:
            match = pattern.match(file)
            if match:
                number = int(match.group(1))
                new_name = re.sub(pattern, f"{number:02d}. ", file)
                old_path = os.path.join(root, file)
                new_path = os.path.join(root, new_name)
                os.rename(old_path, new_path)

        # Rename directories
        for dir in dirs:
            match = pattern.match(dir)
            if match:
                number = int(match.group(1))
                new_name = re.sub(pattern, f"{number:02d}. ", dir)
                old_path = os.path.join(root, dir)
                new_path = os.path.join(root, new_name)
                os.rename(old_path, new_path)

if __name__ == "__main__":
    notes_directory = input("Enter the path to your notes directory: ")
    migrate_files_and_directories(notes_directory)
    print("Migration complete.")