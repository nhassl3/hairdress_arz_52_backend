from pydantic import BaseModel, Field, field_validator
from datetime import date, time, datetime
from uuid import UUID
from typing import Optional


class CreateHairdresserWorkPatterns(BaseModel):

    salon_id: int = Field(gt=0, description="ID салона")
    weekday: int = Field(ge=1, le=7, description="День недели (1=Пн, 7=Вс)")
    shift_start: time = Field(description="Начало смены (HH:MM:SS)")
    shift_end: time = Field(description="Конец смены (HH:MM:SS)")
    effective_from: date = Field(description="Дата начала действия шаблона")
    effective_to: Optional[date] = Field(None, description="Дата окончания (NULL = бессрочно)")

    @field_validator("shift_end")
    def validate_shift_end(cls, v: time, info) -> time:
        shift_start = info.data.get("shift_start")
        if shift_start and v <= shift_start:
            raise ValueError("shift_end must be greater than shift_start")
        return v

    @field_validator("effective_to")
    def validate_effective_to(cls, v: Optional[date], info) -> Optional[date]:
        effective_from = info.data.get("effective_from")
        if v and effective_from and v < effective_from:
            raise ValueError("effective_to must be >= effective_from")
        return v




class UpdateHairdresserWorkPatterns(BaseModel):

    salon_id: Optional[int] = Field(None, gt=0)
    weekday: Optional[int] = Field(None, ge=1, le=7)
    shift_start: Optional[time] = None
    shift_end: Optional[time] = None
    effective_from: Optional[date] = None
    effective_to: Optional[date] = None

    @field_validator("shift_end")
    def validate_shift_end(cls, v: Optional[time], info) -> Optional[time]:
        shift_start = info.data.get("shift_start")
        if v and shift_start and v <= shift_start:
            raise ValueError("shift_end must be greater than shift_start")
        return v



class AdminHairdresserWorkPatterns(CreateHairdresserWorkPatterns):

    id: int
    hairdresser_id: UUID
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True