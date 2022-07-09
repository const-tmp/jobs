from aiogram.dispatcher.filters.state import State, StatesGroup


class MainMenuState(StatesGroup):
    start = State()


class WorkerMenuState(StatesGroup):
    city = State()
    bio = State()
    phone = State()


class HirerMenuState(StatesGroup):
    city = State()
    text = State()
    phone = State()


class ModeratorState(StatesGroup):
    menu = State()
    desc = State()
    skills = State()
    profs = State()
