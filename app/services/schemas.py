
from decimal import Decimal

from datetime import timedelta
from pydantic import BaseModel

from app.bookings.schemas import AdminBookings
from app.hairdresser_services.schemas import HairdresserServices


class CreateServices(BaseModel):

    service_name:str
    duration:timedelta
    price:Decimal
    description:str


class UpdateServices(BaseModel):
    service_name:str| None = None
    duration:timedelta| None = None
    price:Decimal| None = None
    description:str | None = None

class AdminService(CreateServices):
    id:int
    bookings:list["AdminBookings"] | None =None
    hairdresser_services:list["HairdresserServices"] | None =None


