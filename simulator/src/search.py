class Search:
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
