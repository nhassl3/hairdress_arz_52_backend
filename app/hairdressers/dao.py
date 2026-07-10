from app.hairdressers.models import Hairdressers
from app.dao.base import BaseDao
from sqlalchemy.orm import selectinload


class HairdressersDao(BaseDao):
    model = Hairdressers
    _load_options = [
        selectinload(Hairdressers.bookings),
        selectinload(Hairdressers.hairdresser_salons),
        selectinload(Hairdressers.hairdresser_services),
        selectinload(Hairdressers.work_patterns),
        selectinload(Hairdressers.schedules),
        selectinload(Hairdressers.user)
    ]

