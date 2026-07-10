from datetime import datetime
from sqlalchemy import String, DateTime, Text, BigInteger, ForeignKey, func, CheckConstraint, Index
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship
import uuid
from app.database import Base


class Bookings(Base):
    __tablename__ = "bookings"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    username: Mapped[str] = mapped_column(String, ForeignKey("users.username"), nullable=False)
    hairdresser_id: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), ForeignKey("hairdressers.id"), nullable=False)
    service_id: Mapped[int] = mapped_column(ForeignKey("services.id"), nullable=False)
    salon_id: Mapped[int] = mapped_column(ForeignKey("salons.id"), nullable=False)
    starts_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    ends_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    description: Mapped[str] = mapped_column(Text, server_default="")
    status: Mapped[str] = mapped_column(String, nullable=False, server_default="pending")
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now(),
                                                 onupdate=func.now())

    __table_args__ = (
        CheckConstraint("ends_at > starts_at", name="chk_shifts_booking"),
        CheckConstraint("status IN ('pending', 'confirmed', 'completed', 'cancelled', 'no_show')", name="chk_status"),
        Index("ix_bookings_hairdresser_id_starts_at", "hairdresser_id", "starts_at"),
        Index("ix_bookings_salon_id_starts_at", "salon_id", "starts_at"),
        Index("ix_bookings_username_starts_at", "username", "starts_at"),
    )

    # Relationships
    user: Mapped["Users"] = relationship(back_populates="bookings")
    hairdresser: Mapped["Hairdressers"] = relationship(back_populates="bookings")
    service: Mapped["Services"] = relationship(back_populates="bookings")
    salon: Mapped["Salons"] = relationship(back_populates="bookings")