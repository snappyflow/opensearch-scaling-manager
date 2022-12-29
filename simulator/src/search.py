class SearchDescription:
    def __init__(
        self,
        search_type: str,
        heap_load_percent: float,
        cpu_load_percent: int,
        memory_load_percent: float,
    ):
        self.search_type = search_type
        self.heap_load_percent = heap_load_percent
        self.cpu_load_percent = cpu_load_percent
        self.memory_load_percent = memory_load_percent


class SearchState:
    def __init__(
        self,
        position: int,
        time_hh_mm_ss: str,
        search_type: str,
        count: int
    ):
        self.position = position
        self.time_hh_mm_ss = time_hh_mm_ss
        self.search_type = search_type
        self.count = count


class Search:
    def __init__(
        self,
        searches: list[SearchState]
    ):
        self.searches = searches
