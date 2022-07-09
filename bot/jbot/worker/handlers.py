from logging import getLogger

from aiogram import Dispatcher
from aiogram.dispatcher import FSMContext
from aiogram.types import Message, ContentType

from jbot.moder.handlers import new_cv_message, new_moderate_cv_kb, new_cv
from jbot.worker.menu import worker_bio_menu, worker_phone_menu
from jbot.proto.jobs_pb2 import City, SetPhoneRequest, User
from jbot.proto.jobs_pb2_grpc import APIStub
from jbot.states import WorkerMenuState

log = getLogger(__name__)


async def worker_city_handler(msg: Message, state: FSMContext):
    client: APIStub = msg.bot['client']
    ct = await client.GetOrCreateCity(City(name=msg.text))
    log.info(f'{ct = }')
    log.info(f'{ct.name = }')
    await state.update_data({'city': ct.name})
    await worker_bio_menu(msg.bot, msg.chat.id)


async def worker_bio_handler(msg: Message, state: FSMContext):
    await state.update_data({'bio': msg.text})
    await worker_phone_menu(msg.bot, msg.chat.id)


async def worker_phone_handler(msg: Message, state: FSMContext):
    client: APIStub = msg.bot['client']

    res = await client.SetPhone(
        SetPhoneRequest(id=msg.chat.id, phone=msg.contact.phone_number)
    )

    data = await state.get_data()
    cv = await new_cv(client, msg.chat.id, data['city'], data['bio'])

    moderators = await client.GetAllUsers(User(role='user'))
    # moderators = await client.GetAllUsers(User(role='moderator'))

    text = new_cv_message(cv)
    kb = new_moderate_cv_kb(cv.base.id)
    for m in moderators.users:
        await msg.bot.send_message(m.base.id, text, reply_markup=kb)

    # await msg.answer(f'<pre>{html.escape(str(cv))}</pre>')
    await msg.answer(f'Резюме сохранено')


async def worker_phone_text_handler(msg: Message):
    await msg.answer('Нажмите кнопку "Поделиться контактом"')


async def hirer_city_handler(msg: Message, state: FSMContext):
    await state.update_data({'city': msg.text})


def setup_worker_handlers(dp: Dispatcher):
    dp.register_message_handler(
        worker_city_handler,
        state=WorkerMenuState.city
    )
    dp.register_message_handler(
        worker_bio_handler,
        state=WorkerMenuState.bio
    )
    dp.register_message_handler(
        worker_phone_handler,
        state=WorkerMenuState.phone,
        content_types=[ContentType.CONTACT]
    )
    dp.register_message_handler(
        worker_phone_text_handler,
        state=WorkerMenuState.phone,
        content_types=[ContentType.TEXT]
    )
