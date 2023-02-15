from scipy.interpolate import InterpolatedUnivariateSpline
import numpy as np
import math

class Index:
    """
    Represents and holds the data fetched from the configuration for index addition
    """
    def __init__(self,
                index_count:int,
                primary_count:int,
                replica_count:int):
        self.index_count = index_count
        self.primary_count = primary_count
        self.replica_count = replica_count
        
class IndexAddition:
    """
    Parses the configuration and creates relevant index objects,
    Performs aggregation of index addition.
    """
    def __init__(self,
                states: list[dict]
                ):
        self.states = states
        
    def aggregate_index_addition(self,initial_index_count, start_time_hh_mm_ss: str, frequency_minutes:int):
        """
        Produces cumulative index count over time period and returns a list of aggregated index count
        for given duration.

        :param start_time_hh_mm_ss: start time in hh_mm_ss in 24-hour format.
        :param duration_minutes: duration of point generation in minutes
        :param frequency_minutes: gap between the resultant points
        :return: array of int containing resultant index aggregation points
        """
        time_of_day = []
        total_index_count = initial_index_count
        index_addition = []
        start_time_hour = int(start_time_hh_mm_ss.split("_")[0])
        start_time_minutes = int(start_time_hh_mm_ss.split("_")[1])
        if start_time_minutes > 0:
            start_time_hour+=1
        duration_of_day = ((24 - start_time_hour)*60)+ ((60 - start_time_minutes)%60)
        for day in self.states:
            time_of_day.clear()
            index_added = []
            index_count_list = []
            for position in day['pattern']:
                if int(position['time_hh_mm_ss'].split("_")[0]) >= int(start_time_hh_mm_ss.split("_")[0]):
                    time_of_day.append(
                        (int(position['time_hh_mm_ss'].split("_")[0]) - int("0")) * 60
                    )

                index_addition_rate = position.get('index',0)
                if index_addition_rate == 0:
                    index_added.append(Index(0,0,0))
                    index_count_list.append(total_index_count)
                else:
                    index_added.append(Index(index_addition_rate.get('count'),
                                                index_addition_rate.get('primary'),
                                                index_addition_rate.get('replica')
                                                ))
                    total_index_count+= index_addition_rate.get('count')
                    index_count_list.append(total_index_count)

                intervals = int(duration_of_day/frequency_minutes)
                if start_time_hh_mm_ss == "00_00_00":
                    x = np.linspace(0, 24*60, intervals)
                else:
                    start = int(start_time_hh_mm_ss.split("_")[0]) 
                    x = np.linspace(start, duration_of_day, intervals) 

            order = 1
            s = InterpolatedUnivariateSpline(
                time_of_day, index_count_list, k=order
                )

            y = [min(int(math.ceil(max(i, initial_index_count))), total_index_count) for i in s(x)]
            for val in y:
                index_addition.append(val)

        return index_addition
