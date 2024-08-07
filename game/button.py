import pygame


class Button():
    def __init__(self, ai_settings, screen, msg):
        self.screen = screen
        self.screen_rect = screen.get_rect()
        # 设置按钮的属性
        self.width, self.height = 200, 50
        self.button_color = (0, 255, 0)
        self.text_color = (255, 255, 255)
        self.font = pygame.font.SysFont(None, 48)
        # 创建一个按钮的rect对象，并使其居中
        self.rect = pygame.Rect(0, 0, self.width, self.height)
        self.rect.center = self.screen_rect.center

        # 按钮的标签只需创建一次
        self.prep_msg(msg)

    def prep_msg(self, msg):
        """将msg渲染为图像，并使其在按钮上居中"""
        # 将msg渲染为图像
        self.msg_image = self.font.render(msg, True, self.text_color,
                                          self.button_color)
        # 根据文本图像创建一个rect
        self.msg_image_rect = self.msg_image.get_rect()
        # 并将其center属性设置为按钮的center属性
        self.msg_image_rect.center = self.rect.center

    def draw_button(self):
        """绘制一个用颜色填充的按钮，再绘制文本"""
        # 用button_color来绘制按钮的矩形
        self.screen.fill(self.button_color, self.rect)
        # 传递一幅图像以及与该图像相关联的rect对象，从而在屏幕上绘制文本图像
        self.screen.blit(self.msg_image, self.msg_image_rect)
