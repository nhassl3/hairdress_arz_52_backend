from datetime import datetime

from sqlalchemy import String, Boolean,  DateTime, func
from sqlalchemy.orm import mapped_column, relationship, Mapped

from app.database import Base


class Users(Base):
    __tablename__ = 'users'

    username: Mapped[str] = mapped_column(String, primary_key=True)
    full_name: Mapped[str] = mapped_column(String, nullable=False)
    phone_number: Mapped[str] = mapped_column(String, unique=True, nullable=False)\

    is_verified: Mapped[bool] = mapped_column(Boolean, nullable=False, server_default="false")
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now(),
                                                 onupdate=func.now())
    last_login: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)

    # Relationships
    bookings: Mapped[list["Bookings"]] = relationship(back_populates="user")
    admin: Mapped["Admins | None"] = relationship(back_populates="user", cascade="all, delete-orphan")
    hairdresser: Mapped["Hairdressers | None"] = relationship(back_populates="user", cascade="all, delete-orphan")

