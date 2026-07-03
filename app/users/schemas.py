import re
from datetime import datetime

from pydantic import BaseModel, field_validator


class UserRegister(BaseModel):

    username: str
    full_name: str
    phone_number: str

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

class UpdateUser(BaseModel):
    full_name: str | None = None
    phone_number: str | None = None
    is_verified: bool | None = None

    @field_validator("phone_number")
    def validate_phone(cls, v: str | None) -> str | None:
        if v is None:
            return v
        if not re.match(r"^\+7\d{10}$", v):
            raise ValueError("Phone number must be in format: +7XXXXXXXXXX (10 digits after +7)")
        return v


class ReplaceUser(BaseModel):
    full_name: str
    phone_number: str
    is_verified: bool

    @field_validator("phone_number")
    def validate_phone(cls, v: str) -> str:
        if not re.match(r"^\+7\d{10}$", v):
            raise ValueError("Phone number must be in format: +7XXXXXXXXXX (10 digits after +7)")
        return v


class AdminUser(UserRegister):


    is_verified: bool
    created_at: datetime
    updated_at: datetime
    last_login: datetime