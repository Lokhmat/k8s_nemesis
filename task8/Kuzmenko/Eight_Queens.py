# This solution presented by Roman Kuzmenko

from itertools import permutations

class GenerateColumnChoices:
    def generate(self):
        return list(range(8))

class GeneratePermutations:
    def generate(self, choices):
        return permutations(choices)

class CheckConfiguration:
    def check(self, configurations):
        valid_solutions = []
        for config in configurations:
            if self.is_valid(config):
                valid_solutions.append(config)
        return valid_solutions

    def is_valid(self, config):
        n = len(config)
        for i in range(n):
            for j in range(i + 1, n):
                if config[i] == config[j] or abs(config[i] - config[j]) == j - i:
                    return False
        return True

class SaveSolutions:
    def save(self, valid_solutions):
        for solution in valid_solutions:
            print("Solution:", solution)

gen_choices = GenerateColumnChoices()
gen_permutations = GeneratePermutations()
checker = CheckConfiguration()
saver = SaveSolutions()

column_choices = gen_choices.generate()
permutations = gen_permutations.generate(column_choices)
valid_solutions = checker.check(permutations)
print(f"Total solutions: {len(valid_solutions)}")
saver.save(valid_solutions)
