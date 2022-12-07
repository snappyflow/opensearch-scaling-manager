class Search:
    def __init__(
            self,
            search_type: str,
            probability: float,
            cpu_load_percent: int,
            memory_load_percent: float,
    ):
        self.search_type = search_type
        self.probability = probability
        self.cpu_load_percent = cpu_load_percent
        self.memory_load_percent = memory_load_percent
