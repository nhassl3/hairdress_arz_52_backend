from decimal import Decimal
from typing import Optional

from fastapi import APIRouter

from app.exceptions import NotFoundElement, NoFieldsToUpdate, AlreadyExistsElement
from app.services.dao import ServicesDao
from app.services.schemas import AdminService, CreateServices, UpdateServices

from sqlalchemy.exc import IntegrityError


router = APIRouter(
    prefix="/admin",
    tags=["Services"]
)

@router.get('/services', response_model=list[AdminService])
async def get_all_services(skip: int = 0, limit: int = 100):
    return await ServicesDao.find_all(skip=skip, limit=limit)


@router.get('/filter_services/', response_model=AdminService)
async def get_filter_service( id:Optional[int] = None,
                              service_name:Optional[str]=None,
                              duration:Optional[dict]=None,
                              price:Optional[Decimal]=None,
                              description:Optional[str]=None,
                              skip: int = 0, limit: int = 100):
    filters = {}
    if id:
        filters['id'] = id
    if service_name:
        filters['service_name'] = service_name
    if duration:
        filters['duration'] = duration
    if price:
        filters['price'] = price
    if description:
        filters['description'] = description
    if filters:
        services = await ServicesDao.find_by_filter(skip=skip, limit=limit, **filters)
    else:
        services = await ServicesDao.find_all(skip=skip, limit=limit)

    return services


@router.post('/services', response_model=AdminService)
async def create_service(service: CreateServices):
    try:
        new_services = await ServicesDao.add(**service.model_dump())
        return new_services
    except IntegrityError:
        raise AlreadyExistsElement


@router.patch('/services/{service_id}', response_model=AdminService)
async def partial_update_service(service_id: int, service_data: UpdateServices):
    existing_service = await ServicesDao.find_one_or_none(id=service_id)
    if not existing_service:
        raise NotFoundElement

    update_data = service_data.model_dump(exclude_unset=True)
    if not update_data:
        raise NoFieldsToUpdate

    updated_service = await ServicesDao.update(
        filters={"id": service_id},
        data=update_data
    )
    return updated_service