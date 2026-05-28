
from uuid import UUID
from datetime import datetime, date, time

from pydantic import BaseModel, field_validator


class CreateHairdresserSchedules(BaseModel):

    hairdresser_id: UUID
    salon_id: int
    pattern_id: int
    work_date : date
    shift_start : datetime
    shift_end : datetime
    source:str
    comment:str

    @field_validator("shift_end")
    def validate_shift_end(cls, v: time, info) -> time:
        shift_start = info.data.get("shift_start")
        if shift_start and v <= shift_start:
            raise ValueError("shift_end must be greater than shift_start")
        return v


class UpdateHairdresserSchedules(BaseModel):

    shift_start: datetime | None = None
    shift_end: datetime | None = None
    is_available: bool | None = None
    comment: str | None = None

class AdminHairdresserSchedules(CreateHairdresserSchedules):
    id:int
    created_at:datetime
    updated_at:datetime
