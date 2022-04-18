
import os
import sys
import argparse
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt

from plotter import *


CSV_NAMES='query_results.csv'


def main(benchmarks):
    # parse args
    if not benchmarks:
        print('no benchmarks were given')
        exit(1)

    print(benchmarks)
    # TODO: check if outpput path exits
    
    # Loop over the readed files
    for bn in benchmarks:
        print('\nAnalysing benchmark:', bn)

        # Read the csv from that benchmark as a panda
        benchmark_pd = pd.read_csv(bn+'/'+CSV_NAMES)
        print(benchmark_pd)
        
        benchmark = Benchmark(bn, benchmark_pd)
        benchmark.plot_secuence()





if __name__ == "__main__":
    benchmarks = []
    for item in range(1 , len(sys.argv), 1):
        benchmarks.append(sys.argv[item])
    main(benchmarks)