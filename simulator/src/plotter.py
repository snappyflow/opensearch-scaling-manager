import os.path

import matplotlib.pyplot as plt

from config_parser import get_source_code_dir
from constants import SIMULATION_GRAPH_FILE_NAME

font1 = {'family': 'serif', 'color': 'red', 'size': 15}
font2 = {'family': 'serif', 'color': 'darkred', 'size': 10}


def plot_data_points(cluster_objects):
    data_ingestion_over_time = []
    cpu_usage_over_time = []
    mem_usage_over_time = []
    cluster_status_over_time = []
    date_time_points = []
    for cluster_obj in cluster_objects:
        date_time_points.append(cluster_obj.date_time)
        data_ingestion_over_time.append(cluster_obj._ingestion_rate)
        cpu_usage_over_time.append(cluster_obj.cpu_usage_percent)
        mem_usage_over_time.append(cluster_obj.memory_usage_percent)
        cluster_status_over_time.append(cluster_obj.status)
    plt.subplot(4, 1, 1)
    plt.ylabel('Ingestion Rate (in GB/hr)', font2)
    plt.plot(date_time_points, data_ingestion_over_time)

    plt.subplot(4, 1, 2)
    plt.ylabel('Used CPU %', font2)
    plt.plot(date_time_points, cpu_usage_over_time)

    plt.subplot(4, 1, 3)
    plt.ylabel('Used Memory %', font2)
    plt.plot(date_time_points, mem_usage_over_time)

    plt.subplot(4, 1, 4)
    plt.ylabel('Cluster State', font2)
    plt.plot(date_time_points, cluster_status_over_time)

    plt.xlabel('Datetime -->', font2)

    # save the figure
    print('saving graph')
    file_path = os.path.join(get_source_code_dir(), SIMULATION_GRAPH_FILE_NAME)
    plt.savefig(file_path)

    # display the graph
    # plt.show()
