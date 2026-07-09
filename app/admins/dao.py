
from app.admins.models import Admins
from app.dao.base import BaseDao



class AdminsDao(BaseDao):
    model = Admins