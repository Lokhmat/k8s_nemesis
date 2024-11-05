import sys

def main():
    for line in sys.stdin:
        words = line.split()
        n = len(words)
        for i in range(n):
            shifted_line = ' '.join(words[i:] + words[:i])
            print(shifted_line)

if __name__ == "__main__":
    main()