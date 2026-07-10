from datetime import datetime
import uuid

from sqlalchemy import String, Boolean, ForeignKey,  DateTime, func,  UUID, Text
from sqlalchemy.orm import mapped_column, relationship, Mapped

from app.database import Base


class Hairdressers(Base):
    __tablename__ = 'hairdressers'

    id: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    username: Mapped[str] = mapped_column(String, ForeignKey("users.username"), unique=True, nullable=False)
    image_url: Mapped[str] = mapped_column(String, nullable=False )
    is_active: Mapped[bool] = mapped_column(Boolean, nullable=False, server_default="true")
    description: Mapped[str] = mapped_column(Text, nullable=False)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now(),
                                                 onupdate=func.now())

    # Relationships
    user: Mapped["Users"] = relationship(back_populates="hairdresser")
    hairdresser_salons: Mapped[list["HairdresserSalons"]] = relationship(back_populates="hairdresser",
                                                                        cascade="all, delete-orphan")
    hairdresser_services: Mapped[list["HairdresserServices"]] = relationship(back_populates="hairdresser",
                                                                            cascade="all, delete-orphan")
    work_patterns: Mapped[list["HairdresserWorkPatterns"]] = relationship(back_populates="hairdresser",
                                                                         cascade="all, delete-orphan")
    schedules: Mapped[list["HairdresserSchedules"]] = relationship(back_populates="hairdresser",
                                                                  cascade="all, delete-orphan")
    bookings: Mapped[list["Bookings"]] = relationship(back_populates="hairdresser")