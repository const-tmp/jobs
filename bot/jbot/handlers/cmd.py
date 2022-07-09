from logging import getLogger

from aiogram.types import Message
from grpc import RpcError

from jbot.main_menu.menu import main_menu
from jbot.proto import jobs_pb2_grpc, jobs_pb2

# "sender_chat":{"id":-1001744592207,"title":"HVC Orders","type":"channel"}
# "sender_chat":{"id":-1001738223265,"title":"HVC Reports","type":"channel"}
from jbot.states import MainMenuState

log = getLogger(__name__)


async def cmd_start(msg: Message):
    await MainMenuState.start.set()

    client: jobs_pb2_grpc.APIStub = msg.bot['client']

    user = None
    try:
        user = await client.GetUserByID(
            jobs_pb2.IDRequest(id=msg.from_user.id)
        )
    except RpcError as e:
        log.error(e)
        # status = e.code()
        # StatusCode.INTERNAL

    if user is None:
        name = msg.from_user.first_name
        if msg.from_user.last_name:
            name += f' {msg.from_user.last_name}'
        user = await client.CreateUser(
            jobs_pb2.User(
                base=jobs_pb2.Base(id=msg.from_user.id),
                name=name,
                role='user'
            )
        )

    # await msg.answer(user)

    match user.role:
        # case 'user':
        #     pass
        case _:
            await main_menu(msg.bot, msg.chat.id)


async def chtest(msg: Message):
    print(msg)
