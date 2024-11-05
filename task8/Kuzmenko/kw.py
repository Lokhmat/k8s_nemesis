# This solution presented by Roman Kuzmenko

class LineStorage:
    def __init__(self):
        self._lines = []
    
    def add_line(self, line):
        self._lines.append(line)
    
    def get_line(self, index):
        return self._lines[index] if index < len(self._lines) else None
    
    def get_all_lines(self):
        return self._lines.copy()


class CircularShifter:
    def __init__(self, line_storage):
        self._line_storage = line_storage
        self._shifts = []
    
    def generate_shifts(self):
        for line in self._line_storage.get_all_lines():
            words = line.split()
            for i in range(len(words)):
                shift = words[i:] + words[:i]
                self._shifts.append(shift)
    
    def get_shift(self, index):
        return self._shifts[index] if index < len(self._shifts) else None
    
    def get_all_shifts(self):
        return self._shifts.copy()


class Alphabetizer:
    def __init__(self, circular_shifter):
        self._circular_shifter = circular_shifter
        self._sorted_shifts = []
    
    def sort_shifts(self):
        self._sorted_shifts = sorted(
            self._circular_shifter.get_all_shifts(), 
            key=lambda shift: shift[len(shift) // 2].lower()
        )
    
    def get_sorted_shift(self, index):
        return self._sorted_shifts[index] if index < len(self._sorted_shifts) else None
    
    def get_all_sorted_shifts(self):
        return self._sorted_shifts.copy()


class Formatter:
    def __init__(self, alphabetizer):
        self._alphabetizer = alphabetizer 
        self._formatted_output = []
    
    def format_shifts(self):
        for shift in self._alphabetizer.get_all_sorted_shifts():
            center_index = len(shift) // 2
            left_part = " ".join(shift[:center_index])
            keyword = shift[center_index]
            right_part = " ".join(shift[center_index + 1:])
            formatted = f"{left_part} [{keyword}] {right_part}"
            self._formatted_output.append(formatted)
    
    def get_formatted_output(self):
        return self._formatted_output.copy()


class KWICSystem:
    def __init__(self):
        self._line_storage = LineStorage()
        self._circular_shifter = CircularShifter(self._line_storage)
        self._alphabetizer = Alphabetizer(self._circular_shifter)
        self._formatter = Formatter(self._alphabetizer)
    
    def add_line(self, line):
        self._line_storage.add_line(line)
    
    def process(self):
        self._circular_shifter.generate_shifts()
        self._alphabetizer.sort_shifts()
        self._formatter.format_shifts()
    
    def get_output(self):
        return self._formatter.get_formatted_output()


if __name__ == "__main__":
    kwic_system = KWICSystem()
    
    print("Enter lines of text (press Enter twice to finish):")
    input_lines = []
    while True:
        line = input()
        if not line:
            break
        input_lines.append(line)
    
    for line in input_lines:
        kwic_system.add_line(line)
    
    kwic_system.process()
    output = kwic_system.get_output()
    
    for entry in output:
        print(entry)

