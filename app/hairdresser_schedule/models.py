
from datetime import datetime, date
from sqlalchemy import String, Boolean, DateTime, Text, BigInteger, ForeignKey, func, CheckConstraint, Index, Date
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship
import uuid
from app.database import Base

class HairdresserSchedules(Base):
    __tablename__ = "hairdresser_schedules"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    hairdresser_id: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), ForeignKey("hairdressers.id"), nullable=False)
    salon_id: Mapped[int] = mapped_column(ForeignKey("salons.id"), nullable=False)
    pattern_id: Mapped[int | None] = mapped_column(BigInteger, ForeignKey("hairdresser_work_patterns.id"),
                                                   nullable=True)
    work_date: Mapped[date] = mapped_column(Date, nullable=False)
    shift_start: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    shift_end: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    is_available: Mapped[bool] = mapped_column(Boolean, nullable=False, server_default="true")
    source: Mapped[str] = mapped_column(String, nullable=False, server_default="pattern")
    comment: Mapped[str | None] = mapped_column(Text, nullable=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now(),
                                                 onupdate=func.now())

    __table_args__ = (
        CheckConstraint("shift_end > shift_start", name="chk_schedule_time"),
        CheckConstraint("source IN ('pattern', 'manual', 'override')", name="chk_source"),
        Index("ix_hairdresser_schedules_hairdresser_id_work_date", "hairdresser_id", "work_date"),
        Index("ix_hairdresser_schedules_hairdresser_id_shift_start", "hairdresser_id", "shift_start"),
        Index("ix_hairdresser_schedules_salon_id_work_date", "salon_id", "work_date"),
    )

    # Relationships
    hairdresser: Mapped["Hairdressers"] = relationship(back_populates="schedules")
    salon: Mapped["Salons"] = relationship(back_populates="schedules")
    pattern: Mapped["HairdresserWorkPatterns | None"] = relationship(back_populates="schedules")