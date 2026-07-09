from fastapi import FastAPI

# 1. Независимые таблицы
from app.users.models import Users
from app.salons.models import Salons
from app.services.models import Services

# 2. Таблицы, которые ссылаются на независимые
from app.admins.models import Admins
from app.hairdressers.models import Hairdressers

# 3. Таблицы "многие-ко-многим"
from app.hairdresser_salons.models import HairdresserSalons
from app.hairdresser_services.models import HairdresserServices

# 4. Таблицы с внешними ключами на предыдущие
from app.hairdresser_work_patterns.models import HairdresserWorkPatterns
from app.hairdresser_schedule.models import HairdresserSchedules
from app.bookings.models import Bookings
from app.sessions.models import Sessions

from app.users.router import router as users_router
from app.salons.router import router as salons_router
from app.services.router import router as services_router
from app.admins.router import router as admins_router


app = FastAPI()


app.include_router(users_router)
app.include_router(salons_router)
app.include_router(services_router)
app.include_router(admins_router)