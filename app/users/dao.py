
from app.users.models import Users
from app.dao.base import BaseDao



class UsersDao(BaseDao):
    model = Users


