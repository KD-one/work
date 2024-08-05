class GameStats():
    """跟踪游戏的统计信息"""

    def __init__(self, ai_settings):
        """初始化统计信息"""
        self.ai_settings = ai_settings
        self.reset_stats()
        # 游戏刚启动时处于非活动状态
        self.game_active = False
        # 最高得分
        # 重新点击play按钮开启新一轮游戏时，不重置最高得分。除非重新实例化一个Game_Stats 或 重启后端
        self.high_score = 0

    def reset_stats(self):
        """初始化在游戏运行期间可能变化的统计信息"""
        self.ships_left = self.ai_settings.ship_limit
        self.score = 0
        self.level = 1