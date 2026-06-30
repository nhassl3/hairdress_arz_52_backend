from typing import Dict, Any, Optional

from app.database import async_session_maker
from sqlalchemy import select, insert, delete, update



class BaseDao:

    model = None

    @classmethod
    async def find_all(cls):
        async with async_session_maker() as session:
            query = select(cls.model)
            result = await session.execute(query)
            return  result.scalars().all()


    @classmethod
    async def find_by_id(cls, id:int):
        async with async_session_maker() as session:
            query = select(cls.model).where(cls.model.id == id)
            result = await session.execute(query)
            return  result.scalars().first()

    @classmethod
    async def find_one_or_none(cls, **filter):
        async with async_session_maker() as session:
            query=select(cls.model).filter_by(**filter)
            result = await session.execute(query)
            return  result.scalar_one_or_none()

    @classmethod
    async def find_by_filter(cls, skip: int = 0, limit: int = 100, **filter):
        async with async_session_maker() as session:
            query = select(cls.model)
            for key, value in filter.items():
                if hasattr(cls.model, key) and value is not None:
                    query = query.filter(getattr(cls.model, key) == value)

            query = query.offset(skip).limit(limit)  # пагинация в БД
            result = await session.execute(query)
            return result.scalars().all()

    @classmethod
    async def update(cls, filters: Dict[str, Any], data: Dict[str, Any]) -> Optional[model]:
        if not filters:
            raise ValueError("Filters are required for update")

        if not data:
            raise ValueError("No data to update")

        async with async_session_maker() as session:
            query = update(cls.model).filter_by(**filters).values(**data).returning(cls.model)
            result = await session.execute(query)
            await session.commit()
            return result.scalar_one_or_none()

    @classmethod
    async def delete_by_id(cls, id:int):
        async with async_session_maker() as session:
            query = delete(cls.model).where(cls.model.id == id).returning(cls.model)
            await session.execute(query)
            await session.commit()

    @classmethod
    async def delete(cls, **filters):
        async with async_session_maker() as session:
            query = delete(cls.model).filter_by(**filters)
            await session.execute(query)
            await session.commit()


    @classmethod
    async def add(cls, **kwargs):
        async with async_session_maker() as session:
            query = insert(cls.model).values(**kwargs).returning(cls.model)
            result = await session.execute(query)
            await session.commit()
            return  result.scalar_one_or_none()

