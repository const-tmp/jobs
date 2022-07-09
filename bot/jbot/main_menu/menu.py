from logging import getLogger

from aiogram import Bot
from aiogram.types import ReplyKeyboardMarkup, KeyboardButton, Message

from jbot.main_menu.const import find_work, find_workers, main_menu_message

log = getLogger(__name__)

main_menu_kb = ReplyKeyboardMarkup([
    [KeyboardButton(find_work), KeyboardButton(find_workers)]
])


async def main_menu(bot: Bot, chat_id: int) -> Message:
    return await bot.send_message(
        chat_id,
        main_menu_message,
        reply_markup=main_menu_kb
    )
