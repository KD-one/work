class Settings:
    """存储《外星人入侵》的所有设置的类"""

    def __init__(self):
        """初始化游戏的静态设置"""
        # 屏幕设置
        self.screen_width = 1200
        self.screen_height = 800
        self.bg_color = (230, 230, 230)
        # 飞船设置
        self.ship_limit = 3
        # 子弹设置
        self.bullet_width = 3
        self.bullet_height = 15
        self.bullet_color = 60, 60, 60
        # 子弹数量限制表
        self.bullets_allowed = 3
        # 外星人设置
        # fleet_drop_speed表示有外星人撞到屏幕边缘时，外星人群向下移动的速度
        self.fleet_drop_speed = 10
        # 加快游戏节奏的速度
        self.speedup_scale = 1.3
        # 外星人点数的提高速度
        self.score_scale = 1.5
        # 随游戏进行而变化的设置
        self.initialize_dynamic_settings()

    def initialize_dynamic_settings(self):
        """初始化随游戏进行而变化的设置"""
        self.ship_speed_factor = 1
        self.bullet_speed_factor = 1
        self.alien_speed_factor = 0.2

        # fleet_direction为1表示向右；为-1表示向左
        self.fleet_direction = 1

        # 一个外星人值几分（随着游戏的进行，我们将提高每个外星人值的点数）
        self.alien_points = 2

    def increase_speed(self):
        """提高速度设置"""
        self.ship_speed_factor *= self.speedup_scale
        self.bullet_speed_factor *= self.speedup_scale
        self.alien_speed_factor *= self.speedup_scale
        self.alien_points = int(self.alien_points * self.score_scale)