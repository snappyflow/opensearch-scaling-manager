import os.path

import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
from open_search_simulator import timeit

from config_parser import get_source_code_dir
from constants import SIMULATION_GRAPH_FILE_NAME

font1 = {'family': 'serif', 'color': 'red', 'size': 12}
font2 = {'family': 'serif', 'color': 'darkred', 'size': 10}

fig, axes = plt.subplots(1, 7, figsize=(20, 15), dpi=150)


@timeit
def plot_data_points(cluster_objects, skip_data_ingestion=False, skip_search_query=False):
    graph_count = 8
    data_ingestion_over_time = []
    simple_search_query_over_time = []
    medium_search_query_over_time = []
    complex_search_query_over_time = []
    cpu_usage_over_time = []
    mem_usage_over_time = []
    heap_usage_over_time = []
    disk_util_over_time = []
    cluster_status_over_time = []
    nodes_over_time = []
    date_time_points = []
    index_count = []
    for cluster_obj in cluster_objects:
        date_time_points.append(cluster_obj.date_time)
        simple_search_query_over_time.append(cluster_obj._simple_query_rate)
        medium_search_query_over_time.append(cluster_obj._medium_query_rate)
        complex_search_query_over_time.append(cluster_obj._complex_query_rate)
        data_ingestion_over_time.append(cluster_obj._ingestion_rate)
        cpu_usage_over_time.append(cluster_obj.cpu_usage_percent)
        mem_usage_over_time.append(cluster_obj.memory_usage_percent)
        heap_usage_over_time.append(cluster_obj.heap_usage_percent)
        disk_util_over_time.append(cluster_obj.disk_usage_percent)
        cluster_status_over_time.append(cluster_obj.status)
        nodes_over_time.append(cluster_obj.total_nodes_count)
        index_count.append(cluster_obj.instantaneous_index_count)

    if not skip_data_ingestion:
        plt.subplot(graph_count, 1, 1)
        plt.ylabel('Ingestion Rate (in GB/hr)', font2)
        plt.plot(date_time_points, data_ingestion_over_time)
    if not skip_search_query:
        plt.subplot(graph_count, 1, 2)
        plt.ylabel('Search query rate', font2)
        plt.plot(date_time_points, simple_search_query_over_time)
        plt.plot(date_time_points, medium_search_query_over_time)
        plt.plot(date_time_points, complex_search_query_over_time)
        plt.legend(["simple", "medium", "complex"], loc="upper right", prop={'size': 5})

    plt.subplot(graph_count, 1, 3)
    plt.ylabel('Used CPU %', font2)
    plt.plot(date_time_points, cpu_usage_over_time)

    plt.subplot(graph_count, 1, 4)
    plt.ylabel('Used Memory %', font2)
    plt.plot(date_time_points, mem_usage_over_time)

    plt.subplot(graph_count, 1, 5)
    plt.ylabel('Used Heap %', font2)
    plt.plot(date_time_points, heap_usage_over_time)

    plt.subplot(graph_count, 1, 6)
    plt.ylabel('Used Disk %', font2)
    plt.plot(date_time_points, disk_util_over_time)

    # plt.subplot(graph_count, 1, 7)
    # plt.ylabel('Cluster State', font2)
    # plt.plot(date_time_points, cluster_status_over_time)

    plt.subplot(graph_count, 1, 7)
    plt.ylabel('Node Count', font2)
    plt.plot(date_time_points, nodes_over_time)

    plt.subplot(graph_count, 1, 8)
    plt.ylabel('Index Count', font2)
    plt.plot(date_time_points, index_count)

    plt.subplots_adjust(hspace=0.1)
    plt.xlabel('Datetime -->', font2)
    fig.tight_layout()
    plt.gcf().autofmt_xdate()

    # save the figure
    print('saving graph')
    file_path = os.path.join(get_source_code_dir(), SIMULATION_GRAPH_FILE_NAME)
    plt.savefig(file_path)

    # display the graph
    # plt.show()
