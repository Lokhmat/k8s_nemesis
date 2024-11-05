def is_safe_to_place_queen(board, row, col, n):
    for i in range(row):
        if board[i] == col:
            return False

    for i, j in zip(range(row - 1, -1, -1), range(col - 1, -1, -1)):
        if board[i] == j:
            return False

    for i, j in zip(range(row - 1, -1, -1), range(col + 1, n)):
        if board[i] == j:
            return False

    return True

def solve_queens(board, row, n, solutions):
    if row == n:
        solutions.append(board.copy())
        return
    for col in range(n):
        if is_safe_to_place_queen(board, row, col, n):
            board[row] = col
            solve_queens(board, row + 1, n, solutions)

def print_solutions(solutions, n):
    for idx, solution in enumerate(solutions, 1):
        print(f"Solution {idx}:")
        for row in range(n):
            line = ""
            for col in range(n):
                if solution[row] == col:
                    line += "Q "
                else:
                    line += ". "
            print(line)
        print("\n")

def eight_queens(n=8):
    board = [-1] * n
    solutions = []
    solve_queens(board, 0, n, solutions)
    print(f"Found {len(solutions)} solutions for {n} queens.\n")
    print_solutions(solutions, n)

if __name__ == "__main__":
    eight_queens()