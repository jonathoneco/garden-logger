# Index Strategy
## Current Status
I've added support for a settings menu that allows the user to set the indexing strategy for the current directory. Currently, any indexing has been done manually and the ui elements have only been stubbed. Now, I need to implement actual indexing

## Configuration
The indexing strategy for the current directory is defined by a .index file in the directory. It is a json file that describes the indexing strategy. There'll be a key for the current indexing strategy, if there is no .index file then no indexing strategy is applied. Strategy specific configuration defined below

Indexed directories will label any files or directories in them with {index} - {filename} with the dash for easy parsing and visual separation

### Required Functionality
- An enum for the index strategies
- Get .index file for current directory
- Parse .index file
- Write .index file when strategy is selected
- delete .index file when no strategy is selected

## Strategies
These strategies will need to be able to perform certain tasks for indexing to work as expected
For now, we'll only support numeric

- Apply an Indexing strategy to the current directory
- Remove an Indexing strategy from the current directory
- Validate the indexing for the current directory
- Create a new file with the indexing Strategy

## Numeric
This is a simple numeric indexing strategy.

It should take as configuration whether or not to prioritize  directories

> [!INFO] Dir-Priority Indexing
> /
> 1 - Example Dir 1
> 2 - Example Dir 2
> 3 - Example File 1
> 4 - Example File 2

> [!INFO] Not Dir-Priority Indexing
> /
> 1 - Example File 3
> 2 - Example Dir 3
> 3 - Example Dir 4
> 4 - Example File 4

### Functionality
- Should be 1-indexed
- Find Next Index
    - Should increment the highest index
    - If directories are prioritized, should handle that accordingly
- Index Shifting
    - Will need to support index shifting for toggling directory prioritization, and for moving files
- Validation
    - Should validate that every file is properly indexed and there are not gaps in indexes
- Strip index
    - Should be able to strip indexes from files for removing indexing or when moving to another directory
- Fresh Application
    - When applying to a new directory, should set baseline indices naively / alphabetically
    - Not sure if this is possible, but if we could save the most recently applied index to the file's metadata so we can "recover" indexing if it's removed then reapplied

## None
No indexing being done
