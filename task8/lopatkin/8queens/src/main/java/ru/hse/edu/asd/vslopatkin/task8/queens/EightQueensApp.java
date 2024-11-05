package ru.hse.edu.asd.vslopatkin.task8.queens;

import java.util.ArrayList;
import java.util.List;

class Position {
    private int row;
    private int col;

    public Position(int row, int col) {
        this.row = row;
        this.col = col;
    }

    public int getRow() {
        return row;
    }

    public int getCol() {
        return col;
    }

    public boolean isThreatening(Position other) {
        return this.row == other.row || this.col == other.col ||
                Math.abs(this.row - other.row) == Math.abs(this.col - other.col);
    }

    @Override
    public String toString() {
        return "(" + row + ", " + col + ")";
    }
}

class Queen {
    private Position position;

    public Queen(int row, int col) {
        this.position = new Position(row, col);
    }

    public Position getPosition() {
        return position;
    }

    public boolean isThreatening(Queen other) {
        return this.position.isThreatening(other.getPosition());
    }
}

class ChessBoard {
    private int size;
    private List<Queen> queens;

    public ChessBoard(int size) {
        this.size = size;
        this.queens = new ArrayList<>();
    }

    public boolean placeQueen(Queen queen) {
        for (Queen q : queens) {
            if (q.isThreatening(queen)) {
                return false;
            }
        }
        queens.add(queen);
        return true;
    }

    public void removeQueen() {
        if (!queens.isEmpty()) {
            queens.remove(queens.size() - 1);
        }
    }

    public boolean isComplete() {
        return queens.size() == size;
    }

    public List<Queen> getQueens() {
        return new ArrayList<>(queens);
    }

    public int getSize() {
        return size;
    }
}

class Solver {
    private ChessBoard board;
    private List<List<Position>> solutions;

    public Solver(int size) {
        this.board = new ChessBoard(size);
        this.solutions = new ArrayList<>();
    }

    public void solve() {
        placeQueens(0);
    }

    private void placeQueens(int row) {
        int size = board.getSize();
        if (row == size) {
            saveSolution();
            return;
        }

        for (int col = 0; col < size; col++) {
            Queen queen = new Queen(row, col);
            if (board.placeQueen(queen)) {
                placeQueens(row + 1);
                board.removeQueen();
            }
        }
    }

    private void saveSolution() {
        List<Position> solution = new ArrayList<>();
        for (Queen queen : board.getQueens()) {
            solution.add(queen.getPosition());
        }
        solutions.add(solution);
    }

    public List<List<Position>> getSolutions() {
        return solutions;
    }
}

class SolutionCollector {
    public void displaySolutions(List<List<Position>> solutions) {
        for (List<Position> solution : solutions) {
            System.out.println("Solution:");
            for (Position position : solution) {
                System.out.print(position + " ");
            }
            System.out.println("\n");
        }
    }
}

public class EightQueensApp {
    public static void main(String[] args) {
        int size = 8;
        Solver solver = new Solver(size);
        solver.solve();

        SolutionCollector collector = new SolutionCollector();
        collector.displaySolutions(solver.getSolutions());
    }
}