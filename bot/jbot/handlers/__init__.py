from aiogram import Dispatcher

from jbot.handlers.cmd import cmd_start, chtest
from jbot.main_menu.handlers import setup_main_menu_handlers
from jbot.moder.handlers import setup_moder_handlers
from jbot.worker.handlers import setup_worker_handlers


def setup_handlers(dp: Dispatcher):
    dp.register_message_handler(cmd_start, commands=['start'], state='*')
    dp.register_channel_post_handler(chtest)

    setup_main_menu_handlers(dp)
    setup_worker_handlers(dp)
    setup_moder_handlers(dp)
