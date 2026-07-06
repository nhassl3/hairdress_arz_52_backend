
from app.services.models import Services
from app.dao.base import BaseDao
from sqlalchemy.orm import selectinload


class ServicesDao(BaseDao):
    model = Services
    _load_options = [
        selectinload(Services.hairdresser_services),
        selectinload(Services.bookings)
    ]
