from uuid import UUID


from pydantic import BaseModel

class HairdresserSalons(BaseModel):

    hairdresser_id: UUID
    salon_id: int