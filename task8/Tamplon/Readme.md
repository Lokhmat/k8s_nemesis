# Project Overview

This project contains two Python programs: `8queens.py` and `kwic.py`.
for Eight Queens (8Q) and Key Word in Context (KWIC) problems accordingly.

## 8 Queens

The `8queens.py` program solves the classic 8 Queens problem, 
where the goal is to place 8 queens on a chessboard such that no two queens threaten each other.

### How to Run

1. Run the program using the following command:

   ```sh
   python 8queens.py
   
### Example Output
The program will output the number of solutions found and display each solution in a human-readable format. Here is an example of the output:
```
Found 92 solutions for 8 queens.

Solution 1:

Q . . . . . . .   
. . . . Q . . .   
. . . . . . . Q   
. . . Q . . . .   
. . . . . Q . .   
. Q . . . . . .   
. . . . . . Q .   
. . Q . . . . .   

...
```

## KWIC

The `kwic.py` program implements a Key Word In Context (KWIC) index. 
It reads lines of text from standard input, generates all circular shifts of each line, 
sorts the shifts alphabetically, and prints the sorted shifts.  

### How to Run

1. Run the program using the following command:

   ```bash
   python input_module.py text_for_kwic | python circular_shift_module.py | python alphabetizer_module.py | python output_module.py output.txt
   
### Example Input
The input file `text_for_kwic` contains example input:
```
design is hard  
hello world
```
### Example Output
Here is an example of the output:

```
design is hard
hard design is
hello world
is hard design
world hello
```
