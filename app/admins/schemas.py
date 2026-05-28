from uuid import UUID
from datetime import datetime

from pydantic import BaseModel, Field

class CreateAdmins(BaseModel):

    username: str
    level_right:int

class UpdateAdmins(BaseModel):

    level_right: int = Field(ge=1, le=10, description="Новый уровень прав")

class Admins(CreateAdmins):
    id: UUID
    created_at: datetime
    updated_at: datetime