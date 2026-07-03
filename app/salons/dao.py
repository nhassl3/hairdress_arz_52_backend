

from app.salons.models import Salons
from app.dao.base import BaseDao



class SalonDAO(BaseDao):
    model = Salons