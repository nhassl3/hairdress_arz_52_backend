from sqlalchemy.orm import selectinload

from app.admins.models import Admins
from app.dao.base import BaseDao



class AdminsDao(BaseDao):
    model = Admins
    _load_options = [
        selectinload(Admins.user)
    ]