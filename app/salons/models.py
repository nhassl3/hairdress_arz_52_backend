from datetime import datetime

from sqlalchemy import String, Boolean,  Integer,  DateTime, func, Index
from sqlalchemy.orm import mapped_column, relationship, Mapped

from app.database import Base



class Salons(Base):
    __tablename__ = 'salons'
    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    salon_name: Mapped[str] = mapped_column(String, nullable=False)
    address: Mapped[str] = mapped_column(String, nullable=False)
    phone: Mapped[str] = mapped_column(String, nullable=False)
    is_active: Mapped[bool] = mapped_column(Boolean, nullable=False, server_default="true")
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now(),
                                                 onupdate=func.now())

    __table_args__ = (
        Index("ix_salon_salon_name", "salon_name"),
    )

    # Relationships
    hairdresser_salons: Mapped[list["HairdresserSalons"]] = relationship(back_populates="salon",
                                                                        cascade="all, delete-orphan")
    work_patterns: Mapped[list["HairdresserWorkPatterns"]] = relationship(back_populates="salon")
    schedules: Mapped[list["HairdresserSchedules"]] = relationship(back_populates="salon")
    bookings: Mapped[list["Bookings"]] = relationship(back_populates="salon")