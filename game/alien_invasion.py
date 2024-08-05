import pygame

from button import Button
from game_stats import GameStats
from scoreboard import Scoreboard
from settings import Settings
from ship import Ship
import game_functions as gf


def run_game():
    # 初始化游戏并创建一个屏幕对象
    pygame.init()
    # 创建一个Settings实例
    ai_settings = Settings()
    screen = pygame.display.set_mode((ai_settings.screen_width, ai_settings.screen_height))
    pygame.display.set_caption("飞船大战")

    # 创建play按钮
    play_button = Button(ai_settings, screen, "Play")

    # 创建一艘飞船
    ship = Ship(ai_settings, screen)
    # 创建一个用于存储子弹的编组
    bullets = pygame.sprite.Group()
    # 创建一个用于存储外星人的编组
    aliens = pygame.sprite.Group()

    # 更新外星人位置
    gf.create_fleet(ai_settings, screen, aliens, ship)

    # 创建一个用于存储游戏统计信息的实例
    stats = GameStats(ai_settings)
    # 创建记分牌
    sb = Scoreboard(ai_settings, screen, stats)

    # 开始游戏主循环
    while True:
        # 监视键盘和鼠标事件
        gf.check_events(ai_settings, screen, stats, play_button, ship, aliens, bullets, sb)
        # 如果游戏处于活动状态，就更新飞船的位置，子弹的位置，外星人的位置
        if stats.game_active:
            # 更新飞船的位置
            ship.update()
            # 更新子弹的位置
            gf.update_bullets(ai_settings, screen, bullets, aliens, ship, stats, sb)
            # 更新外星人的位置
            gf.update_aliens(ai_settings, screen, stats, ship, aliens, bullets, sb)
        # 使用更新后的位置来刷新屏幕
        gf.update_screen(ai_settings, screen, ship, aliens, bullets, stats, play_button, sb)


run_game()
