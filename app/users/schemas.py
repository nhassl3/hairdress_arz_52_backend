import re
import uuid
from datetime import datetime

from pydantic import BaseModel, field_validator

from app.bookings.schemas import AdminBookings
from app.hairdressers.schemas import AdminHairdresser
from app.admins.schemas import Admins as AdminOut


class UserRegister(BaseModel):

    username: str
    full_name: str
    phone_number: str
    email: str

    @field_validator("username")
    def validate_username(cls, v: str) -> str:
        if not re.match(r"^[a-zA-Z0-9_-]+$", v):
            raise ValueError("Username can only contain letters, numbers, underscore and hyphen")
        return v.lower()

    @field_validator("phone_number")
    def validate_phone(cls, v: str) -> str:
        if not re.match(r"^\+7\d{10}$", v):
            raise ValueError("Phone number must be in format: +7XXXXXXXXXX (10 digits after +7)")
        return v

    @field_validator("email")
    def validate_email(cls, v: str) -> str:
        if not re.match(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$", v):
            raise ValueError("Invalid email format")
        return v.lower()


class UpdateUser(BaseModel):
    full_name: str | None = None
    phone_number: str | None = None
    is_verified: bool | None = None
    email: str | None = None
    role: str | None = None

    @field_validator("phone_number")
    def validate_phone(cls, v: str | None) -> str | None:
        if v is None:
            return v
        if not re.match(r"^\+7\d{10}$", v):
            raise ValueError("Phone number must be in format: +7XXXXXXXXXX (10 digits after +7)")
        return v

    @field_validator("email")
    def validate_email(cls, v: str | None) -> str | None:
        if v is None:
            return v
        if not re.match(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$", v):
            raise ValueError("Invalid email format")
        return v.lower()


class AdminUser(UserRegister):

    uid: uuid.UUID
    is_verified: bool
    role: str
    created_at: datetime
    updated_at: datetime
    last_login: datetime
    bookings: list[AdminBookings] = []
    hairdresser: AdminHairdresser | None = None
    admin: AdminOut | None = None
