from time import sleep

import pygame

from alien import Alien
from bullet import Bullet
from button import Button


def check_keydown_events(ai_settings, screen, event, ship, bullets):
    """响应按键"""
    if event.key == pygame.K_RIGHT:
        ship.moving_right = True
    elif event.key == pygame.K_LEFT:
        ship.moving_left = True
    elif event.key == pygame.K_SPACE:
        fire_bullet(ai_settings, screen, ship, bullets)


def check_keyup_events(event, ship):
    """响应松开"""
    if event.key == pygame.K_RIGHT:
        ship.moving_right = False
    elif event.key == pygame.K_LEFT:
        ship.moving_left = False


def check_play_button(ai_settings, screen, stats, play_button, ship, aliens, bullets, mouse_x, mouse_y, sb):
    """在玩家单击Play按钮时开始新游戏"""
    button_clicked = play_button.rect.collidepoint(mouse_x, mouse_y)
    if button_clicked and not stats.game_active:
        # 重置游戏设置
        ai_settings.initialize_dynamic_settings()
        # 隐藏光标
        pygame.mouse.set_visible(False)
        # 重置游戏统计信息
        stats.reset_stats()
        stats.game_active = True

        # 重置记分牌图像
        sb.prep_score()
        sb.prep_high_score()
        sb.prep_level()
        sb.prep_ships()

        # 清空外星人列表和子弹列表
        aliens.empty()
        bullets.empty()

        # 创建一群新的外星人，并让飞船居中
        create_fleet(ai_settings, screen, aliens, ship)
        ship.center_ship()


def check_events(ai_settings, screen, stats, play_button, ship, aliens, bullets, sb):
    """响应按键和鼠标事件"""
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            pygame.quit()
            return
        elif event.type == pygame.KEYDOWN:
            check_keydown_events(ai_settings, screen, event, ship, bullets)
        elif event.type == pygame.KEYUP:
            check_keyup_events(event, ship)
        elif event.type == pygame.MOUSEBUTTONDOWN:
            mouse_x, mouse_y = pygame.mouse.get_pos()
            check_play_button(ai_settings, screen, stats, play_button, ship, aliens, bullets, mouse_x, mouse_y, sb)


def update_screen(ai_settings, screen, ship, aliens, bullets, stats, play_button, sb):
    """
    更新屏幕上的图像
    """
    # 每次循环时都重绘屏幕
    screen.fill(ai_settings.bg_color)
    # 绘制飞船
    ship.blitme()
    # 绘制外星人
    aliens.draw(screen)
    # 绘制子弹
    for bullet in bullets.sprites():
        bullet.draw_bullet()
    # 显示得分
    sb.show_score()
    # 绘制其他所有游戏元素后再绘制这个按钮是为了为了让Play按钮位于其他所有屏幕元素上面
    # 如果游戏处于非活动状态，就绘制Play按钮
    if not stats.game_active:
        play_button.draw_button()
    # 让最近绘制的屏幕可见
    pygame.display.flip()


def update_bullets(ai_settings, screen, bullets, aliens, ship, stats, sb):
    """更新子弹的位置，并删除已消失的子弹, 并检查是否有子弹击中了外星人"""
    # 对编组调用update()时，编组将自动对其中的每个精灵调用update()，
    # 因此代码行bullets.update()将为编组bullets中的每颗子弹调用bullet.update()
    bullets.update()

    # 删除已消失的子弹
    for bullet in bullets.copy():
        # 检查子弹是否到达屏幕顶部
        if bullet.rect.bottom <= 0:
            bullets.remove(bullet)

    check_bullet_alien_collisions(ai_settings, screen, ship, aliens, bullets, stats, sb)


def check_high_score(stats, sb):
    """检查是否诞生了新的最高得分"""
    if stats.score > stats.high_score:
        stats.high_score = stats.score
        sb.prep_high_score()


def check_bullet_alien_collisions(ai_settings, screen, ship, aliens, bullets, stats, sb):
    """响应子弹和外星人的碰撞"""

    # 删除发生碰撞的子弹和外星人
    # 与外星人碰撞的子弹都是字典collisions中的一个键；而与每颗子弹相关的值都是一个列表，其中包含该子弹撞到的外星人
    collisions = pygame.sprite.groupcollide(bullets, aliens, True, True)

    if collisions:
        # 每个值aliens都是一个列表，包含被同一颗子弹击中的所有外星人
        for aliens in collisions.values():
            # 更新分数
            stats.score += ai_settings.alien_points * len(aliens)
            # 更新计分牌
            sb.prep_score()
            # 检查是否诞生了新的最高得分
            check_high_score(stats, sb)

    # 检查是否已删除所有外星人
    if len(aliens) == 0:
        # 删除现有的所有子弹
        bullets.empty()
        # 加快游戏节奏（增加子弹、飞船、外星人的速度）
        ai_settings.increase_speed()
        # 提高等级
        stats.level += 1
        # 重新绘制等级图像
        sb.prep_level()
        # 创建一个新的外星人群
        create_fleet(ai_settings, screen, aliens, ship)


def fire_bullet(ai_settings, screen, ship, bullets):
    """如果还没有达到限制，就发射一颗子弹"""
    # 创建新子弹，并将其加入到编组bullets中
    if len(bullets) < ai_settings.bullets_allowed:
        new_bullet = Bullet(ai_settings, screen, ship)
        bullets.add(new_bullet)


def get_number_aliens_x(ai_settings, alien_width):
    """计算每行可容纳多少个外星人"""
    available_space_x = ai_settings.screen_width - 2 * alien_width
    number_aliens_x = int(available_space_x / (2 * alien_width))
    return number_aliens_x


def get_number_rows(ai_settings, ship_height, alien_height):
    """计算屏幕可容纳多少行外星人"""
    available_space_y = ai_settings.screen_height - (3 * alien_height) - ship_height
    number_rows = int(available_space_y / (2 * alien_height))
    return number_rows


def create_alien(ai_settings, screen, aliens, alien_number, row_number):
    """创建一个外星人并将其放在当前行"""
    alien = Alien(ai_settings, screen)
    alien.x = alien.rect.width + 2 * alien.rect.width * alien_number
    alien.rect.x = alien.x
    alien.y = alien.rect.height + 2 * alien.rect.height * row_number
    alien.rect.y = alien.y
    aliens.add(alien)


def create_fleet(ai_settings, screen, aliens, ship):
    """创建外星人群"""
    # 创建一个外星人，并计算每行可容纳多少个外星人
    alien = Alien(ai_settings, screen)
    number_aliens_x = get_number_aliens_x(ai_settings, alien.rect.width)
    number_rows = get_number_rows(ai_settings, ship.rect.height, alien.rect.height)
    # 创建外星人群
    for row_number in range(number_rows):
        for alien_number in range(number_aliens_x):
            create_alien(ai_settings, screen, aliens, alien_number, row_number)


def change_fleet_direction(ai_settings, aliens):
    """将整群外星人下移，并改变它们的方向"""
    for alien in aliens.sprites():
        alien.rect.y += ai_settings.fleet_drop_speed
    ai_settings.fleet_direction *= -1


def check_fleet_edges(ai_settings, aliens):
    """有外星人到达边缘时采取相应的措施"""
    for alien in aliens.sprites():
        if alien.check_edges():
            change_fleet_direction(ai_settings, aliens)
            break


def check_aliens_bottom(ai_settings, stats, screen, ship, aliens, bullets, sb):
    """有外星人到达屏幕低端时采取相应的措施"""
    screen_rect = screen.get_rect()
    for alien in aliens.sprites():
        if alien.rect.bottom >= screen_rect.bottom:
            # 像飞船被撞到一样进行处理
            ship_hit(ai_settings, stats, screen, ship, aliens, bullets, sb)
            break


def update_aliens(ai_settings, screen, stats, ship, aliens, bullets, sb):
    """
    检查是否有外星人位于屏幕边缘，并更新整群外星人的位置
    """
    check_fleet_edges(ai_settings, aliens)
    aliens.update()

    # 检测外星人和飞船之间的碰撞
    if pygame.sprite.spritecollideany(ship, aliens):
        ship_hit(ai_settings, stats, screen, ship, aliens, bullets, sb)

    # 检查是否有外星人到达屏幕低端
    check_aliens_bottom(ai_settings, stats, screen, ship, aliens, bullets, sb)


def ship_hit(ai_settings, stats, screen, ship, aliens, bullets, sb):
    """响应被外星人撞到的飞船"""
    if stats.ships_left > 0:
        # 将ships_left减1
        stats.ships_left -= 1
        # 更新飞船计分牌(左上角的几条命)
        sb.prep_ships()
        # 清空外星人列表和子弹列表
        aliens.empty()
        bullets.empty()
        # 创建一群新的外星人，并将飞船放到屏幕底端中央
        create_fleet(ai_settings, screen, aliens, ship)
        ship.center_ship()
        # 暂停
        sleep(2)
    else:
        stats.game_active = False
        # 显示光标
        pygame.mouse.set_visible(True)
