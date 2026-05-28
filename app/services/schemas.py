
from decimal import Decimal

from pydantic import BaseModel


class CreateServices(BaseModel):

    service_name:str
    duration:dict
    price:Decimal
    description:str


class UpdateServices(BaseModel):
    service_name:str| None = None
    duration:dict| None = None
    price:Decimal| None = None
    description:str | None = None

class AdminService(CreateServices):
    id:int


