from uuid import UUID


from pydantic import BaseModel

class HairdresserServices(BaseModel):

    hairdresser_id:UUID
    service_id:int
