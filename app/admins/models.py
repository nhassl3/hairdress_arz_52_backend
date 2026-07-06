from datetime import datetime
import uuid

from sqlalchemy import String, Boolean, Column, ForeignKey, Integer, Date, DateTime, func, Index, BigInteger, UUID, \
    Numeric, SmallInteger
from sqlalchemy.orm import mapped_column, Mapped, relationship

from app.database import Base

class Admins(Base):
    __tablename__ = 'admins'

    id: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    username: Mapped[str] = mapped_column(String, ForeignKey("users.username"), unique=True, nullable=False)
    level_right: Mapped[SmallInteger] = mapped_column(Numeric, nullable=False, server_default="1")
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now(),
                                                 onupdate=func.now())

    # Relationships
    user: Mapped["Users"] = relationship(back_populates="admin")