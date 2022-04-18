
import os
import pandas as pd
from datetime import datetime
import matplotlib.pyplot as plt
import matplotlib.dates as mdates


PLOT_FOLDER_NAME='plots'

class Benchmark():

    def __init__(self, folder, pd_obj):
        self.folder = folder
        self.pd_obj = pd_obj
        
        self.create_plot_folder()
        
        self.dates, self.resp_times = self.get_dist_arrays()

    def create_plot_folder(self):
        try:
            # Create Plot folder inside benchmark folder
            os.mkdir(self.folder+'/'+PLOT_FOLDER_NAME) 
        except Exception as e:
            pass

    def plot_secuence(self):

        print('generating graphs')
        # Plotting secuence for the benchmark analysis
        plt.plot(self.dates, self.resp_times)

        plt.show()
    
    #def plot(self, x_arrays, y_arrays, opts):


    def get_dist_arrays(self):
        date_array = []
        resp_time_array = [] 
        for idx, row in self.pd_obj.iterrows():
            # Response time
            resp_time = parse_time_into_milli_secs(row['response time'])
            resp_time_array.append(resp_time)
            #print(row['response time'], resp_time)

            # Date
            date = parse_req_time(row['request time'])
            date_array.append(date)
            #print(row['request time'], date)
        return date_array, resp_time_array
            

def parse_time_into_milli_secs(org_time):
    milli_sec = 0.0
    if 'µs' in org_time:
        org_time = org_time.replace('µs', '')
        milli_sec = float(org_time) / 1000
    elif 'ms' in org_time:
        org_time = org_time.replace('ms', '')
        milli_sec = float(org_time)  
    elif 's' in org_time:
        org_time = org_time.replace('s', '')
        milli_sec = float(org_time) * 1000

    return milli_sec

def parse_req_time(org_time):
    org_time = org_time.split(' +')[0]
    org_time = org_time[:25]
    date = datetime.strptime(org_time, '%Y-%m-%d %H:%M:%S.%f')

    return date



    