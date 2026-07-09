import uuid
from typing import Optional

from fastapi import APIRouter
from sqlalchemy.exc import IntegrityError

from app.exceptions import NotFoundElement, NoFieldsToUpdate, AlreadyExistsElement, HairdresserHasBookings
from app.hairdressers.dao import HairdressersDao
from app.hairdressers.schemas import AdminHairdresser, CreateHairdresser, UpdateHairdresser


router = APIRouter(
    prefix="/admin",
    tags=["Hairdressers"]
)


@router.get('/hairdressers', response_model=list[AdminHairdresser])
async def get_all_hairdressers(skip: int = 0, limit: int = 100):
    return await HairdressersDao.find_all(skip=skip, limit=limit)


@router.get('/hairdresser_filter/', response_model=list[AdminHairdresser])
async def get_filter_hairdressers(
    username: Optional[str] = None,
    is_active: Optional[bool] = None,
    skip: int = 0,
    limit: int = 100
):
    filters = {}
    if username:
        filters['username'] = username
    if is_active is not None:
        filters['is_active'] = is_active

    if filters:
        return await HairdressersDao.find_by_filter(skip=skip, limit=limit, **filters)
    return await HairdressersDao.find_all(skip=skip, limit=limit)


@router.get('/hairdressers/{hairdresser_id}', response_model=AdminHairdresser)
async def get_hairdresser(hairdresser_id: uuid.UUID):
    hairdresser = await HairdressersDao.find_one_or_none(id=hairdresser_id)
    if not hairdresser:
        raise NotFoundElement
    return hairdresser


@router.post('/hairdressers', response_model=AdminHairdresser, status_code=201)
async def create_hairdresser(data: CreateHairdresser):
    try:
        new_hairdresser = await HairdressersDao.add(**data.model_dump())
        return new_hairdresser
    except IntegrityError:
        raise NotFoundElement


@router.patch('/hairdressers/{hairdresser_id}/', response_model=AdminHairdresser)
async def partial_update_hairdresser(hairdresser_id: uuid.UUID, data: UpdateHairdresser):
    existing = await HairdressersDao.find_one_or_none(id=hairdresser_id)
    if not existing:
        raise NotFoundElement

    update_data = data.model_dump(exclude_unset=True)
    if not update_data:
        raise NoFieldsToUpdate

    if "username" in update_data:
        username_user = await HairdressersDao.find_one_or_none(username=update_data["username"])
        if username_user and username_user.id != hairdresser_id:
            raise AlreadyExistsElement

    try:
        updated = await HairdressersDao.update(
            filters={"id": hairdresser_id},
            data=update_data
        )
        return updated
    except IntegrityError:
        raise AlreadyExistsElement


@router.delete('/hairdressers/{hairdresser_id}')
async def delete_hairdresser(hairdresser_id: uuid.UUID):
    existing = await HairdressersDao.find_one_or_none(id=hairdresser_id)
    if not existing:
        raise NotFoundElement
    try:
        await HairdressersDao.delete(id=hairdresser_id)
        return {"detail": "Hairdresser deleted"}
    except IntegrityError:
        raise HairdresserHasBookings

