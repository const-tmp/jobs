from logging import getLogger

from aiogram import Bot
from aiogram.types import Message, ReplyKeyboardMarkup, KeyboardButton, \
    ReplyKeyboardRemove

from jbot.worker.const import city_menu_message, new_city_menu_message, \
    send_worker_bio_message, share_contact_msg
from jbot.proto.jobs_pb2 import EmptyRequest
from jbot.proto.jobs_pb2_grpc import APIStub
from jbot.states import WorkerMenuState

log = getLogger(__name__)


async def worker_city_menu(bot: Bot, chat_id: int) -> Message:
    await WorkerMenuState.city.set()

    client: APIStub = bot['client']
    cts = await client.GetAllCities(EmptyRequest())
    log.info(f'{cts = }')
    log.info(f'{cts.cities = }')

    if not cts.cities:
        return await bot.send_message(
            chat_id,
            new_city_menu_message,
            reply_markup=ReplyKeyboardRemove()
        )

    kb = ReplyKeyboardMarkup(one_time_keyboard=True, resize_keyboard=True)
    for city in cts.cities:
        kb.add(KeyboardButton(city.name))

    return await bot.send_message(
        chat_id,
        city_menu_message,
        reply_markup=kb
    )


async def worker_bio_menu(bot: Bot, chat_id: int) -> Message:
    await WorkerMenuState.bio.set()
    return await bot.send_message(
        chat_id,
        send_worker_bio_message,
        reply_markup=ReplyKeyboardRemove()
    )


async def worker_phone_menu(bot: Bot, chat_id: int) -> Message:
    await WorkerMenuState.phone.set()
    kb = ReplyKeyboardMarkup()
    kb.add(KeyboardButton('Поделиться номером', request_contact=True))
    return await bot.send_message(
        chat_id,
        share_contact_msg,
        reply_markup=kb
    )
