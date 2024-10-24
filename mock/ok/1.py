
import pygetwindow
import pyautogui
import time

# 获取指定窗口标题
target_window_title = 'Async Updating'

# 获取指定窗口
window = pygetwindow.getWindowsWithTitle(target_window_title)

if window:
# 获取窗口对象
window = window[0]

# 重复执行最小化和还原操作
for i in range(10000):
# 最小化窗口
window.minimize()
time.sleep(0.1)

# 还原窗口
window.restore()
time.sleep(0.1)
