import os
import re
from logging import getLogger

from aiogram import Dispatcher, Bot
from aiogram.dispatcher import FSMContext
from aiogram.types import Message, CallbackQuery, ContentType, \
    InlineKeyboardMarkup, \
    InlineKeyboardButton
from aiogram.utils.callback_data import CallbackData

from jbot.proto.jobs_pb2 import City, CV, IDRequest, Profession, Skill
from jbot.proto.jobs_pb2_grpc import APIStub
from jbot.states import ModeratorState
from jbot.utils import user_link

log = getLogger(__name__)

words_by_comma = re.compile(r',\s*')

CV_CHANNEL_ID = os.environ.get('CV_CHANNEL_ID')
AD_CHANNEL_ID = os.environ.get('AD_CHANNEL_ID')


def new_cv_message(cv: CV) -> str:
    return f'''<b>Новое резюме</b>

<b>Контакт:</b> {user_link(cv.user.base.id, cv.user.name)}

<b>Город:</b> {cv.city.name}

<b>Телефон:</b> {cv.user.phone}

<b>Резюме:</b> {cv.raw}'''


async def new_cv(client: APIStub,
                 user_id: int,
                 city: str,
                 bio: str) -> CV:
    ct = await client.GetOrCreateCity(City(name=city))
    res = await client.CreateCV(CV(
        user_id=user_id,
        city_id=ct.base.id,
        raw=bio
    ))
    log.info(f'{res = }')
    return res


def new_moderate_cv_kb(cvid: int) -> InlineKeyboardMarkup:
    kb = InlineKeyboardMarkup()
    kb.add(InlineKeyboardButton(
        'Модерировать',
        callback_data=cv_cb.new(id=cvid, action='new')
    ))
    return kb


async def get_cv_from_callback(bot: Bot, callback_data: str) -> CV:
    cvid = int(callback_data.split(':')[1])
    client: APIStub = bot['client']
    cv = await client.GetCVByID(IDRequest(id=cvid))
    log.info(f'{cv = }')
    return cv


async def moderate_cv_handler(cb: CallbackQuery, state: FSMContext):
    await cb.message.edit_reply_markup()
    await ModeratorState.menu.set()

    cv = await get_cv_from_callback(cb.bot, cb.data)
    log.info(f'{cv = }')
    await state.update_data({'current_cv': cv})
    data = await state.get_data()

    await edit_cv_menu(
        cb.message.bot,
        cb.message.chat.id,
        data.get('current_cv'),
        data.get('edited_desc'),
        data.get('edited_profs'),
        data.get('edited_skills'),
    )


cv_cb = CallbackData('cv', 'id', 'action')  # cv:<id>:<action>


def moderate_cv_menu_kb(cv: CV) -> InlineKeyboardMarkup:
    kb = InlineKeyboardMarkup(
        row_width=2,
        inline_keyboard=[
            [InlineKeyboardButton(
                'Описание',
                callback_data=cv_cb.new(id=cv.base.id, action='edit_desc')
            )],
            [InlineKeyboardButton(
                'Навыки',
                callback_data=cv_cb.new(id=cv.base.id, action='edit_skills')
            )],
            [InlineKeyboardButton(
                'Профессии',
                callback_data=cv_cb.new(id=cv.base.id, action='edit_profs')
            )],
            [InlineKeyboardButton(
                'Запостить в канал',
                callback_data=cv_cb.new(id=cv.base.id, action='post')
            )],
        ]
    )
    kb.add()
    return kb


async def hirer_city_handler(msg: Message, state: FSMContext):
    await state.update_data({'city': msg.text})


async def cv_edit_desc_handler(cb: CallbackQuery, state: FSMContext):
    await cb.message.edit_text('Пришлите отредактированную версию')
    await ModeratorState.desc.set()
    await load_cv(cb, state)


async def cv_edited_desc_handler(msg: Message, state: FSMContext):
    await state.update_data({'edited_desc': msg.text})
    await send_edit_cv_menu(msg, state)


async def cv_edit_skills_handler(cb: CallbackQuery, state: FSMContext):
    await cb.message.edit_text('Пришлите навыки через запятую')
    await ModeratorState.skills.set()
    await load_cv(cb, state)


async def cv_edited_skills_handler(msg: Message, state: FSMContext):
    skills = []
    for w in words_by_comma.split(msg.text):
        prepared = w.strip().lower().capitalize()
        if prepared and prepared not in skills:
            skills.append(prepared)

    await state.update_data({'edited_skills': skills})
    await send_edit_cv_menu(msg, state)


async def send_edit_cv_menu(msg: Message, state: FSMContext):
    data = await state.get_data()
    await edit_cv_menu(
        msg.bot,
        msg.chat.id,
        data.get('current_cv'),
        data.get('edited_desc'),
        data.get('edited_profs'),
        data.get('edited_skills'),
    )


async def load_cv(cb: CallbackQuery, state: FSMContext):
    data = await state.get_data()
    if not data.get('current_cv'):
        cv = await get_cv_from_callback(cb.bot, cb.data)
        log.info(f'{cv = }')
        await state.update_data({'current_cv': cv})


async def cv_edit_profs_handler(cb: CallbackQuery, state: FSMContext):
    await cb.message.edit_text('Пришлите профессии через запятую')
    await ModeratorState.profs.set()
    await load_cv(cb, state)


async def cv_edited_profs_handler(msg: Message, state: FSMContext):
    profs = []
    for w in words_by_comma.split(msg.text):
        prepared = w.strip().lower().capitalize()
        if prepared and prepared not in profs:
            profs.append(prepared)

    await state.update_data({'edited_profs': profs})

    await send_edit_cv_menu(msg, state)


async def edit_cv_menu(bot: Bot, chat_id: int, cv: CV,
                       edited_desc=None,
                       edited_profs=None,
                       edited_skills=None) -> Message:
    await ModeratorState.menu.set()

    old_skills = ', '.join([s.name for s in cv.skills]) \
        if cv.skills is not None else None
    old_profs = ', '.join([s.name for s in cv.professions]) \
        if cv.professions is not None else None
    edited_skills = ', '.join([s for s in edited_skills]) \
        if edited_skills is not None else None
    edited_profs = ', '.join([s for s in edited_profs]) \
        if edited_profs is not None else None

    text = f'''<b>Редактирование</b>

<b>ID</b>: {cv.base.id}
<b>Город</b>: {cv.city.name}
<b>Контакт</b>: {user_link(cv.user.base.id, cv.user.name)}
<b>Телефон</b>: {cv.user.phone}
<b>Исходное реезюме</b>:
{cv.raw}

<b>Отредактированное реезюме</b>:
{edited_desc if edited_desc else "-"}

<b>Исходные навыки</b>:
{old_skills if old_skills else "-"}


<b>Исходные профессии</b>:
{old_profs if old_profs else "-"}


<b>Отредактированное навыки</b>:
{edited_skills if edited_skills else "-"}


<b>Отредактированное профессии</b>:
{edited_profs if edited_profs else "-"}
'''

    return await bot.send_message(
        chat_id,
        text,
        reply_markup=moderate_cv_menu_kb(cv)
    )


async def save_and_post(cb: CallbackQuery, state: FSMContext):
    await load_cv(cb, state)
    data = await state.get_data()
    cv = data.get('current_cv')
    desc = data.get('edited_desc')
    profs = data.get('edited_profs')
    skills = data.get('edited_skills')
    cv.desc = desc
    del cv.professions[:]
    for p in profs:
        cv.professions.append(Profession(name=p))
    del cv.skills[:]
    for s in skills:
        cv.skills.append(Skill(name=s))

    cv.moderated = True
    cv.moderated_by = cb.message.chat.id

    client: APIStub = cb.bot['client']
    cv = await client.UpdateCV(cv)

    text = cv_message(cv)
    await cb.bot.send_message(CV_CHANNEL_ID, text)
    await cb.message.answer(text)


def cv_message(cv) -> str:
    skills = ' '.join([
        f'#{s.name.replace(" ", "_")}' for s in cv.skills
    ]) if cv.skills is not None else None

    profs = ' '.join([
        f'#{s.name.replace(" ", "_")}' for s in cv.professions
    ]) if cv.professions is not None else None

    return f'''#{cv.city.name}

<b>Навыки</b>: {skills if skills else "-"}

<b>Профессии</b>: {profs if profs else "-"}

<b>Контакт</b>: {user_link(cv.user.base.id, cv.user.name)}
<b>Телефон</b>: {cv.user.phone}
<b>Резюме</b>:
{cv.desc}
'''


def setup_moder_handlers(dp: Dispatcher):
    dp.register_callback_query_handler(
        moderate_cv_handler,
        cv_cb.filter(action='new'),
        state='*'
    )

    dp.register_callback_query_handler(
        cv_edit_desc_handler,
        cv_cb.filter(action='edit_desc'),
        state='*'
    )
    dp.register_message_handler(
        cv_edited_desc_handler,
        state=ModeratorState.desc,
        content_types=ContentType.TEXT
    )
    dp.register_callback_query_handler(
        cv_edit_skills_handler,
        cv_cb.filter(action='edit_skills'),
        state='*'
    )
    dp.register_message_handler(
        cv_edited_skills_handler,
        state=ModeratorState.skills,
        content_types=ContentType.TEXT
    )
    dp.register_callback_query_handler(
        cv_edit_profs_handler,
        cv_cb.filter(action='edit_profs'),
        state='*'
    )
    dp.register_message_handler(
        cv_edited_profs_handler,
        state=ModeratorState.profs,
        content_types=ContentType.TEXT
    )
    dp.register_callback_query_handler(
        save_and_post,
        cv_cb.filter(action='post'),
        state='*'
    )
