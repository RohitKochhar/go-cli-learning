from random import randint

# Quick script to generate 1000 files to run benchmark tests on
for i in range(0, 1000):
    with open(f"./testdata/benchmark/file{i}.csv", "a+") as f:
        f.write("Col1,Col2,Col3\n")
        for j in range(0, 2500):
            f.write(f"Data{j},{randint(0, 10000)},{randint(0, 10000)}\n")
            
