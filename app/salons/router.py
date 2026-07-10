
from fastapi import APIRouter
from sqlalchemy.exc import IntegrityError

from app.exceptions import AlreadyExistsElement, NoFieldsToUpdate, NotFoundElement, SalonHasBookings
from app.salons.dao import SalonDAO
from app.salons.schemas import AdminSalon, CreateSalon, UpdateSalon

router = APIRouter(
    prefix="/admin",
    tags=["Salons"]
)

@router.get('/salons', response_model=list[AdminSalon])
async def get_all_salons(skip: int = 0, limit: int = 10):
    return await SalonDAO.find_all(skip=skip, limit=limit)


@router.post("/salons", response_model=AdminSalon)
async def create_salon(salon:CreateSalon):
    try:
        new_salon = await SalonDAO.add(**salon.model_dump())
        return new_salon
    except IntegrityError:
        raise AlreadyExistsElement



@router.patch("/salons/{salon_id}", response_model=AdminSalon)
async def partial_update_salon(salon_id: int, update_data: UpdateSalon):

    existing_salon = await SalonDAO.find_by_id(salon_id)
    if not existing_salon:
        raise NotFoundElement

    salon_data = update_data.model_dump(exclude_unset=True)
    if not salon_data:
        raise NoFieldsToUpdate

    if "salon_name" in salon_data:
        existing_with_new_name = await SalonDAO.find_one_or_none(salon_name=salon_data["salon_name"])
        if existing_with_new_name and existing_with_new_name.id != salon_id:
            raise AlreadyExistsElement
    updated_salon = await SalonDAO.update(filters={"id": salon_id}, data=salon_data)
    return updated_salon

@router.delete("/salons/{salon_id}")
async def delete_salon(salon_id: int):
    existing_salon = await SalonDAO.find_by_id(salon_id)
    if not existing_salon:
        raise NotFoundElement
    try:
        await SalonDAO.delete_by_id(salon_id)
        return {"detail": "Salon deleted"}
    except IntegrityError:
        raise SalonHasBookings