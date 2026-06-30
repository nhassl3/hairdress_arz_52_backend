from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlalchemy.orm import sessionmaker, DeclarativeBase

from app.config import settings

engine = create_async_engine(settings.async_database_url, echo=True)

async_session_maker = sessionmaker(engine,  class_=AsyncSession, expire_on_commit=False)

class Base(DeclarativeBase):
    pass



from app.users.models import Users
from app.salons.models import Salons
from app.services.models import Services
from app.admins.models import Admins
from app.hairdressers.models import Hairdressers
from app.hairdresser_salons.models import HairdresserSalons
from app.hairdresser_services.models import HairdresserServices
from app.hairdresser_work_patterns.models import HairdresserWorkPatterns
from app.hairdresser_schedule.models import HairdresserSchedules
from app.bookings.models import Bookings