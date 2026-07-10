
from app.users.models import Users
from app.dao.base import BaseDao
from sqlalchemy.orm import selectinload


class UsersDao(BaseDao):
    model = Users
    _load_options = [
        selectinload(Users.bookings),
        selectinload(Users.hairdresser),
        selectinload(Users.admin),
        selectinload(Users.sessions)
    ]


