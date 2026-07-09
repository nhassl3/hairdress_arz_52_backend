import uuid
from datetime import datetime

from pydantic import BaseModel

class CreateHairdresser(BaseModel):

    username: str
    image_url: str
    description: str



class UpdateHairdresser(BaseModel):
    username: str| None = None
    image_url: str|None = None
    description: str|None = None


class AdminHairdresser(BaseModel):
    id: uuid.UUID
    username: str
    image_url: str
    created_at: datetime
    updated_at: datetime


