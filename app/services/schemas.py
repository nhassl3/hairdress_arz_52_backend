
from decimal import Decimal

from asyncpg.pgproto.pgproto import timedelta
from pydantic import BaseModel


class CreateServices(BaseModel):

    service_name:str
    duration:timedelta
    price:Decimal
    description:str


class UpdateServices(BaseModel):
    service_name:str| None = None
    duration:timedelta| None = None
    price:Decimal| None = None
    description:str | None = None

class AdminService(CreateServices):
    id:int


