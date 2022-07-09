import os
from logging import getLogger
from logging.config import dictConfig

import click
import dotenv
import grpc
from aiogram import Bot, Dispatcher
from aiogram.contrib.fsm_storage.memory import MemoryStorage
from aiogram.types import ParseMode
from aiogram.utils import executor

from jbot.handlers import setup_handlers
from jbot.log.logger import setup_logger
from jbot.proto.jobs_pb2_grpc import APIStub

logger = getLogger(__name__)


@click.command()
def cli():
    setup_logger()
    dotenv.load_dotenv()
    token = os.getenv('TG_BOT_TOKEN')
    api_addr = os.getenv('API_ADDR')
    if token is None:
        raise ValueError(f'{token = }')

    bot = Bot(token, parse_mode=ParseMode.HTML)
    dp = Dispatcher(bot, storage=MemoryStorage())

    setup_handlers(dp)

    channel = grpc.aio.insecure_channel(api_addr)
    client = APIStub(channel)
    bot.data.update({'client': client})

    executor.start_polling(dp, skip_updates=False)
