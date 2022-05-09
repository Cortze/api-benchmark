import matplotlib.pyplot as plt
from pathlib import Path


COLORS = {
    'prysm':        ['#ce1e16', '#e68e8a', '#7b120d', '#9fc0b6', '#538e7b', '#204136'], # red
    'lighthouse':   ['#3232ff', '#b2b2ff', '#00007f', '#ffbc8f', '#ff791f', '#994812'], # blue
    'teku':         ['#ffa500', '#ffd27f', '#7f5200', '#bacce4', '#7e99bd', '#465569'], # orange
    'nimbus':       ['#008000', '#99cc99', '#004000', '#d9c5b2', '#a08a74', '#50453a'], # green
    'lodestar':     ['#8c198c', '#cc99cc', '#4c004c', '#9adf7b', '#58a835', '#2c541a'], # purple
    'grandine':     ['#999900', '#cccc7f', '#4c4c00', '#e69da5', '#c82236', '#610510'], # yellow / gold
}


class Options():
    
    def __init__(self): 
        self.att = {
            'benchmark_name': '',
            'fig_size': (10,6),
            'title': '',
            'title_size': 14,
            'x_label': '',
            'y_label': '',
            'label_size': 12,
            'legend_position': '',
            'save_path': '', 
            'marker': ['+'],
            'marker_linestyle': 'None',
            'marker_size': 1,
            'marker_color': ['tab:blue'],
            'legend_label': [''],
        }

class SingleFigure():

    def __init__(self, opts, x_array, y_array):
        self.fig= plt.figure() # opts['fig_size']
        self.ax = self.fig.add_subplot(111)
        self.opts = opts.att
        self.x_array = x_array
        self.y_array = y_array

    def generate_single_plot(self):
        # generate plot
        for idx, _ in enumerate(self.x_array):
            self.ax.plot(self.x_array[idx], self.y_array[idx],
            marker=self.opts['marker'][idx], 
            linestyle=self.opts['marker_linestyle'],
            markersize=self.opts['marker_size'], 
            color=self.opts['marker_color'][idx], 
            label=self.opts['legend_label'][idx])

        # Set the labels
        self.ax.grid(which='major', axis='x', linestyle='--')
        self.ax.set_xlabel(self.opts['x_label'], fontsize =self.opts['label_size'])
        
        self.ax.set_ylabel(self.opts['y_label'], fontsize =self.opts['label_size'])
        #self.ax.set_ylim(bottom=0)
        
        # compose title
        plt.title(self.opts['title'], fontsize = self.opts['title_size'])
        plt.tight_layout()

    def generate_box_plot(self):
        # generate plot
        self.ax.boxplot(self.y_array, labels=self.opts['legend_label'])

        # Set the labels
        self.ax.grid(which='major', axis='x', linestyle='--')
        self.ax.set_xlabel(self.opts['x_label'], fontsize =self.opts['label_size'])
        
        self.ax.set_ylabel(self.opts['y_label'], fontsize =self.opts['label_size'])
        #self.ax.set_ylim(bottom=0)
        
        # compose title
        plt.title(self.opts['title'], fontsize = self.opts['title_size'])
        plt.tight_layout()

    def generate_pie_plot(self):

        # generate plot
        self.ax.pie(self.y_array, explode=(0, 0.1), labels=self.opts['legend_label'], autopct='%1.1f%%',
            shadow=True, startangle=90)
        
        self.ax.axis('equal') 

        # compose title
        plt.title(self.opts['title'], fontsize = self.opts['title_size'])
        plt.tight_layout()

    def generate_bar_plot(self):
        for idx, y in enumerate(self.y_array):
            print(self.x_array[idx], y)
            self.ax.bar(self.x_array[idx], y[0], color=self.opts['marker_color'][0])
            self.ax.bar(self.x_array[idx], y[1], bottom=y[0], color=self.opts['marker_color'][1])

        # Set the labels
        self.ax.grid(which='major', axis='y', linestyle='--')
        
        self.ax.set_ylabel(self.opts['y_label'], fontsize =self.opts['label_size'])
        #self.ax.set_ylim(bottom=0)
        
        # compose title
        plt.title(self.opts['title'], fontsize = self.opts['title_size'])

    def save_to_file(self, outfile):
        figurePath =  Path(__file__).parent / outfile
        plt.savefig(figurePath)

    def show(self):
        plt.show()