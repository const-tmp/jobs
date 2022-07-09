def user_link(chat_id: int, name: str) -> str:
    return f'<a href="tg://user?id={chat_id}">{name}</a>'
