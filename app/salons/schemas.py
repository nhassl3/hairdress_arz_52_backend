from datetime import datetime

from pydantic import BaseModel

class CreateSalon(BaseModel):
    salon_name: str
    address: str
    phone: str
    is_active: bool

class UpdateSalon(BaseModel):
    salon_name: str | None = None
    address: str | None = None
    phone: str | None = None
    is_active: bool | None = None

class AdminSalon(CreateSalon):
    id:int
    created_at: datetime
    updated_at: datetime
