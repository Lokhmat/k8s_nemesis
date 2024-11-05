import sys

def main():
    lines = [line.strip() for line in sys.stdin if line != ""]
    sorted_lines = sorted(lines)
    for line in sorted_lines:
        print(line)

if __name__ == "__main__":
    main()