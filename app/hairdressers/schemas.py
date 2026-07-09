import uuid
from datetime import datetime

from pydantic import BaseModel

from app.bookings.schemas import AdminBookings
from app.hairdresser_salons.schemas import HairdresserSalons
from app.hairdresser_services.schemas import HairdresserServices
from app.hairdresser_work_patterns.schemas import AdminHairdresserWorkPatterns
from app.hairdresser_schedule.schemas import AdminHairdresserSchedules


class HairdresserUser(BaseModel):
    username: str
    full_name: str
    phone_number: str
    email: str

    class Config:
        from_attributes = True


class CreateHairdresser(BaseModel):

    username: str
    image_url: str
    description: str

    class Config:
        from_attributes = True


class UpdateHairdresser(BaseModel):
    username: str | None = None
    image_url: str | None = None
    description: str | None = None
    is_active: bool | None = None


class AdminHairdresser(CreateHairdresser):
    id: uuid.UUID
    is_active: bool
    created_at: datetime
    updated_at: datetime
    user: HairdresserUser | None = None
    hairdresser_salons: list[HairdresserSalons] = []
    hairdresser_services: list[HairdresserServices] = []
    work_patterns: list[AdminHairdresserWorkPatterns] = []
    schedules: list[AdminHairdresserSchedules] = []
    bookings: list[AdminBookings] = []


