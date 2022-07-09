from aiogram import Dispatcher
from aiogram.types import Message

from jbot.main_menu.const import find_work, find_workers
from jbot.worker.menu import worker_city_menu
from jbot.states import MainMenuState, HirerMenuState


async def worker_handler(msg: Message):
    await worker_city_menu(msg.bot, msg.chat.id)


async def hirer_handler(msg: Message):
    await HirerMenuState.city.set()
    await worker_city_menu(msg.bot, msg.chat.id)


def setup_main_menu_handlers(dp: Dispatcher):
    dp.register_message_handler(
        worker_handler,
        text=find_work,
        state=MainMenuState.start
    )
    dp.register_message_handler(
        worker_handler,
        text=find_workers,
        state=MainMenuState.start
    )
