from uuid import UUID
from datetime import datetime


from pydantic import BaseModel, Field

class CreateBookings(BaseModel):

    hairdresser_id: UUID = Field(description="ID мастера")
    service_id: int = Field(gt=0, description="ID услуги")
    salon_id: int = Field(gt=0, description="ID салона")
    starts_at: datetime = Field(description="Дата и время начала")
    description: str = Field(default="", max_length=500, description="Пожелания клиента")



class UpdateBookings(BaseModel):

    username: str | None = None
    hairdresser_id: UUID | None = None
    service_id: int | None = None
    salon_id: int | None = None
    starts_at: datetime | None = None
    status: str | None = None
    description: str | None = None

class AdminBookings(CreateBookings):

    id:int
    username: str
    ends_at: datetime
    status: str
    created_at: datetime
    updated_at: datetime

