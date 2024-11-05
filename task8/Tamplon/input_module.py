import sys

def main():
    if len(sys.argv) < 2:
        print("Usage: python input_module.py <input_filename>", file=sys.stderr)
        sys.exit(1)

    input_filename = sys.argv[1]

    try:
        with open(input_filename, 'r') as file:
            for line in file:
                line = line.strip()
                if line:
                    print(line)
    except FileNotFoundError:
        print(f"File {input_filename} not found.", file=sys.stderr)


if __name__ == "__main__":
    main()