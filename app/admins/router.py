from fastapi import APIRouter

from app.admins.schemas import CreateAdmins, UpdateAdmins
from app.exceptions import NotFoundElement, NoFieldsToUpdate
from app.admins.dao import AdminsDao

from sqlalchemy.exc import IntegrityError


router = APIRouter(
    prefix="/admin",
    tags=["Admins"]
)


@router.get('/admins', response_model=list[CreateAdmins])
async def get_all_admins(skip: int = 0, limit: int = 100):
    return await AdminsDao.find_all(skip=skip, limit=limit)


@router.post('/admins', response_model=CreateAdmins)
async def create_admin(admin: CreateAdmins):
    try:
        new_admin = await AdminsDao.add(**admin.model_dump())
        return new_admin
    except IntegrityError:
        raise NotFoundElement


@router.patch('/admins/{username}', response_model=CreateAdmins)
async def update_admin(admin: UpdateAdmins, username: str):
    existing_admin = await AdminsDao.find_one_or_none(username=username)
    if not existing_admin:
        raise NotFoundElement

    update_data = admin.model_dump(exclude_unset=True)
    if not update_data:
        raise NoFieldsToUpdate

    updated_admin = await AdminsDao.update(
        filters={"username": username},
        data=update_data
    )
    return updated_admin


@router.delete('/admins/{username}')
async def delete_admin(username: str):
    existing_admin = await AdminsDao.find_one_or_none(username=username)
    if not existing_admin:
        raise NotFoundElement

    await AdminsDao.delete(username=username)
    return {"detail": "Admin deleted"}
