
import os
import sys
import argparse
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt

from benchmark import *


CSV_NAMES='query_results.csv'


def main(export_folder, bn_projects):
    # parse args
    if not bn_projects:
        print('no bn_projects were given')
        exit(1)

    print(f"exporting plots to folder: {export_folder}")
    print(bn_projects)

    benchmarks = []
    
    labels = []
    dist_comp = []
    success_ratios = []
    benchmark_names = []
    client_array = []
    set_array = []
    routines_array = []

    # Loop over the readed files
    # Generate sigle plots
    # and compose array of benchmarks with datasets
    for bn in bn_projects:
        print('\nAnalysing benchmark:', bn)
        client, set, routines = parse_benchmark_name(bn)

        client_array.append(client+'_'+routines)
        set_array.append(set)
        routines_array.append(routines)

        # Read the csv from that benchmark as a panda
        benchmark_pd = pd.read_csv(bn+'/'+CSV_NAMES)
        benchmark_names.append(bn)
        
        b = Benchmark(export_folder, bn, benchmark_pd, client, set, routines)
        
        benchmarks.append(b)
        labels.append(client)
        dist_comp.append(b.resp_times[0])
        success_ratios.append(b.ratios['success_rate'])

        b.plot_secuence()


    # compose summary plots from aggregation of clients

    # ------ Success Ratio Comparison ---------
    opts = Options()
    opts.att['title'] = 'Comparison of Success Ratio'
    opts.att['legend_position'] =  ''
    opts.att['marker_color'] =  ['tab:blue', 'tab:red']
    opts.att['legend_label'] =  ['SUCCEED', 'FAILED']
    opts.att['benchmark_name'] =  client_array 

    suc_array = []
    for item in success_ratios:
        suc_row = []
        suc_row.append(item)
        suc_row.append(100.0-item)
        suc_array.append(suc_row)
    
    print(len(client_array), client_array)
    print(len(suc_array), suc_array)

    fig = SingleFigure(opts, client_array, suc_array)

    fig.generate_bar_plot()
    fig_name = Path(export_folder+"/all_comparison_success_ratio.png")
    fig.save_to_file(fig_name)

    # ------ Percentile Plots for the client ---------
    opts = Options()
    opts.att['title'] = 'Response Time Distribution'
    opts.att['y_label'] =  'Response Times (seconds)'
    colors = []
    for item in client_array:
        colors.append(COLORS[item.split('_')[0]][1]) 
    opts.att['marker_color'] =  colors
    opts.att['legend_label'] =  client_array

    fig2 = SingleFigure(opts, [], dist_comp)

    fig2.generate_box_plot()
    fig2_name = Path(export_folder+"/all_comparison_resp_time_percentiles.png")
    print(fig2_name)
    fig2.save_to_file(fig2_name)


    
    
def parse_benchmark_name(folder):
    bn_name = folder.split('/')[-1]

    name = bn_name.split('-')[0]
    meta = name.split('_')
    client = meta[0]
    routines = meta[-1]
    sets = meta[-2] 

    return client, sets, routines





if __name__ == "__main__":
    bn_projects = []
    export_folder = sys.argv[1]
    for item in range(2 , len(sys.argv), 1):
        bn_projects.append(sys.argv[item])
    main(export_folder, bn_projects)