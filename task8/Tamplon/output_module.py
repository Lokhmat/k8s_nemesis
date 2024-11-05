import sys


def main():
    if len(sys.argv) < 2:
        print("Usage: python output_module.py <output_filename>", file=sys.stderr)
        sys.exit(1)

    output_filename = sys.argv[1]

    try:
        with open(output_filename, 'w') as file:
            for line in sys.stdin:
                file.write(line)
    except IOError as e:
        print(f"Error writing to file {output_filename}: {e}", file=sys.stderr)


if __name__ == "__main__":
    main()