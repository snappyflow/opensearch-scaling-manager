import math

import numpy as np
from scipy.interpolate import InterpolatedUnivariateSpline


class SearchStat:
    def __init__(
        self,
        heap_load_percent: float,
        cpu_load_percent: int,
        memory_load_percent: float,
    ):
        self.heap_load_percent = heap_load_percent
        self.cpu_load_percent = cpu_load_percent
        self.memory_load_percent = memory_load_percent


class SearchDescription:
    def __init__(
        self,
        search_type: str,
        search_stat: SearchStat
    ):
        self.search_type = search_type
        self.search_stat = search_stat
        print(self.search_stat)
        print(self.search_type)

    # def __getitem__(self, item):
    #     if self.search_type == item:



class SearchState:
    def __init__(
        self,
        position: int,
        time_hh_mm_ss: str,
        searches: dict
    ):
        self.position = position
        self.time_hh_mm_ss = time_hh_mm_ss
        self.searches = searches


class Search:
    def __init__(
        self,
        searches: list[SearchState]
    ):
        self.searches = searches

    def __str__(self):
        for search in self.searches:
            print(search.time_hh_mm_ss)
            print(search.position)
            print(search.searches)

    # def random_aggregation_points(self, duration_minutes: int, frequency_minutes: int):
    #     print("aggregating random patterns")
    #     # positions to inter/extrapolate
    #     intervals = int(duration_minutes / frequency_minutes)
    #     random_set = np.random.randint(0, self.randomness_percentage, size=intervals)
    #     x = np.linspace(0, duration_minutes, intervals)
    #     return x, random_set

    def data_aggregation_points(
            self, start_time_hh_mm_ss: str, duration_minutes: int, frequency_minutes: int
    ):
        """
        Produce cumulative data points of all events and return an array of resultant aggregation
        :param start_time_hh_mm_ss: start time in hh_mm_ss in 24-hour format, eg. '080000'
        :param duration_minutes: duration of point generation in minutes
        :param frequency_minutes: gap between the resultant points
        :return: array of float containing resultant data aggregation points
        """

        # fits
        time_of_day = []
        search_simple_count = []
        search_medium_count = []
        search_complex_count = []
        search_count = {}
        y = {}

        for search in self.searches:
            time_of_day.append(search.time_hh_mm_ss)
            search_simple_count.append(search.searches.get("simple", 0))
            search_medium_count.append(search.searches.get("medium", 0))
            search_complex_count.append(search.searches.get("complex", 0))

        search_count["simple"] = search_simple_count
        search_count["medium"] = search_medium_count
        search_count["complex"] = search_complex_count

        time_of_day = [int(i.split("_")[0]) * 60 for i in time_of_day]

        print("time_of_day", time_of_day)

        # add missing value of 0th hour
        if time_of_day[0] != 0:
            time_of_day.insert(0, 0)

        # positions to inter/extrapolate
        intervals = int(duration_minutes / frequency_minutes)

        x = np.linspace(0, duration_minutes, intervals)
        # spline order: 1 linear, 2 quadratic, 3 cubic ...
        order = 1
        # do inter/extrapolation
        for search_type in ("simple", "medium", "complex"):
            s = InterpolatedUnivariateSpline(
                time_of_day, search_count[search_type], k=order
            )
            y[search_type] = [math.ceil(i) for i in s(x)]
        return x, y

    def aggregate_data(
            self, start_time_hh_mm_ss: str, duration_minutes: int, frequency_minutes: int
    ):
        x, y1 = self.data_aggregation_points(
            start_time_hh_mm_ss, duration_minutes, frequency_minutes
        )
        #x, y2 = self.random_aggregation_points(duration_minutes, frequency_minutes)
        return x, y1
