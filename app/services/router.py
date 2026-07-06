from typing import Optional

from fastapi import APIRouter

from app.exceptions import NotFoundElement, NoFieldsToUpdate, AlreadyExistsElement, UserHasBookings
from app.users.dao import UsersDao
from app.users.schemas import AdminUser, UserRegister, UpdateUser, ReplaceUser

from sqlalchemy.exc import IntegrityError


router = APIRouter(
    prefix="/admin",
    tags=["Services"]
)

