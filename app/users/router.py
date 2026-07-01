from typing import Optional

from fastapi import APIRouter, HTTPException

from app.exceptions import NotFoundElement, NoFieldsToUpdate, AlreadyExistsElement, UserHasBookings
from app.users.dao import UsersDao
from app.users.schemas import AdminUser, UserRegister, UpdateUser

from sqlalchemy.exc import IntegrityError


router = APIRouter(
    prefix="/admin",
    tags=["Users"]
)


@router.get('/users', response_model=list[AdminUser])
async def get_all_users(skip: int = 0, limit: int = 100):
    return await UsersDao.find_all(skip=skip, limit=limit)


@router.get('/user_filter/', response_model=list[AdminUser])
async def get_filter_users( username:Optional[str] = None,
                          full_name: Optional[str] = None,
                          phone_number:Optional[str] = None,
                          is_verified: Optional[bool] = None,skip: int = 0,limit: int = 100):
    filters = {}
    if username:
        filters['username'] = username
    if full_name:
        filters['full_name'] = full_name
    if phone_number:
        filters['phone_number'] = phone_number
    if is_verified is not None:
        filters['is_verified'] = is_verified

    if filters:
        users = await UsersDao.find_by_filter(skip=skip, limit=limit, **filters)
    else:
        users = await UsersDao.find_all(skip=skip, limit=limit)

    return users



# @router.get('/users/search/', response_model=AdminUser)
# async def search_one_element_user(
#         field: str,
#         value: str,
# ):
#     allowed_fields = {'username', 'full_name', 'phone_number'}
#
#     if field not in allowed_fields:
#         raise NotFoundElement
#
#     user = await UsersDao.find_one_or_none(**{field: value})
#
#     if not user:
#         raise NotFoundElement
#
#     return user


@router.post('/users', response_model=AdminUser, status_code=201)
async def create_user(user: UserRegister):

    try:
        new_user = await UsersDao.add(**user.model_dump())
        return new_user
    except IntegrityError:
        raise AlreadyExistsElement





@router.put('/users/{username}/', response_model=AdminUser)
async def update_user(username: str, user_data: UpdateUser):

    existing_user = await UsersDao.find_one_or_none(username=username)
    if not existing_user:
        raise NotFoundElement

    update_data = user_data.model_dump(exclude_unset=True)
    if not update_data:
        raise NoFieldsToUpdate

    if "phone_number" in update_data:
        phone_user = await UsersDao.find_one_or_none(phone_number=update_data["phone_number"])
        if phone_user and phone_user.username != username:
            raise AlreadyExistsElement


    updated_user = await UsersDao.update(
        filters={"username": username},
        data=update_data
    )
    return updated_user


@router.patch('/users/{username}/', response_model=AdminUser)
async def partial_update_user(username: str, user_data: UpdateUser):

    existing_user = await UsersDao.find_one_or_none(username=username)
    if not existing_user:
        raise NotFoundElement

    update_data = user_data.model_dump(exclude_unset=True)
    if not update_data:
        raise NoFieldsToUpdate

    if "phone_number" in update_data:
        phone_user = await UsersDao.find_one_or_none(phone_number=update_data["phone_number"])
        if phone_user and phone_user.username != username:
            raise AlreadyExistsElement

    updated_user = await UsersDao.update(
        filters={"username": username},
        data=update_data
    )
    return updated_user


@router.delete('/users/{username}/')
async def delete_user(username: str):
    existing_user = await UsersDao.find_one_or_none(username=username)
    if not existing_user:
        raise NotFoundElement
    try:
        await UsersDao.delete(username=username)
        return {"detail": "User deleted"}
    except IntegrityError:
        raise UserHasBookings

