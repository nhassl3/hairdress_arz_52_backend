from sqlalchemy.ext.asyncio import AsyncSession

from app.database import async_session_maker
from sqlalchemy import select, insert, delete, update
from fastapi import APIRouter, Depends, HTTPException

from app.users.dao import UsersDao

router = APIRouter(
    prefix="/",
    tags=[""]
)


@router.get('/users')
async def get_all_users():
    users = await UsersDao.find_all()
    return  users


@router.get('/users/{id}')
async def get_one_user(id:int):
    user = await UsersDao.find_by_id(id)
    return user
