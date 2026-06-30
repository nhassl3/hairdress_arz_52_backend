from typing import Optional

from app.users.models import Users

from app.database import async_session_maker
from sqlalchemy import select, insert, delete, update
from app.dao.base import BaseDao
from app.users.schemas import AdminUser


class UsersDao(BaseDao):
    model = Users


