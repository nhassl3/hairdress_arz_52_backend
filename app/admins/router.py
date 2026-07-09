from decimal import Decimal
from typing import Optional

from fastapi import APIRouter

from app.admins.models import Admins
from app.exceptions import NotFoundElement, NoFieldsToUpdate, AlreadyExistsElement
from app.admins.dao import AdminsDao


from sqlalchemy.exc import IntegrityError


router = APIRouter(
    prefix="/admin",
    tags=["Admins"]
)

@router.get('/admins', response_model=list[Admins])
async def get_all_admins(skip: int = 0, limit: int = 100):
    return await AdminsDao.find_all(skip=skip, limit=limit)

