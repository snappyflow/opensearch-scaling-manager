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
            if int(search.time_hh_mm_ss.split("_")[0]) >= int(start_time_hh_mm_ss.split("_")[0]):
                time_of_day.append(
                    (int(search.time_hh_mm_ss.split("_")[0]) - int(start_time_hh_mm_ss.split("_")[0])) * 60
                )
                search_simple_count.append(search.searches.get("simple", 0))
                search_medium_count.append(search.searches.get("medium", 0))
                search_complex_count.append(search.searches.get("complex", 0))

        search_count["simple"] = search_simple_count
        search_count["medium"] = search_medium_count
        search_count["complex"] = search_complex_count

        # add missing value of 0th hour
        if start_time_hh_mm_ss == "00_00_00" and time_of_day[0] != 0:
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
            y[search_type] = [max(math.ceil(i),0) for i in s(x)]
        return x, y

    def aggregate_data(
            self, start_time_hh_mm_ss: str, duration_minutes: int, frequency_minutes: int
    ):
        x, y1 = self.data_aggregation_points(
            start_time_hh_mm_ss, duration_minutes, frequency_minutes
        )
        #x, y2 = self.random_aggregation_points(duration_minutes, frequency_minutes)
        return x, y1
