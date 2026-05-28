from datetime import date, time, datetime
from sqlalchemy import SmallInteger, Time, Date, DateTime, ForeignKey, func, CheckConstraint, Index
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship
import uuid


from app.database import Base


class HairdresserWorkPatterns(Base):
    __tablename__ = 'hairdresser_work_patterns'

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    hairdresser_id: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), ForeignKey("hairdressers.id"), nullable=False)
    salon_id: Mapped[int] = mapped_column(ForeignKey("salons.id"), nullable=False)
    weekday: Mapped[int] = mapped_column(SmallInteger, nullable=False)
    shift_start: Mapped[time] = mapped_column(Time, nullable=False)
    shift_end: Mapped[time] = mapped_column(Time, nullable=False)
    effective_from: Mapped[date] = mapped_column(Date, nullable=False)
    effective_to: Mapped[date | None] = mapped_column(Date, nullable=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now())
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now(),
                                                 onupdate=func.now())

    __table_args__ = (
        CheckConstraint("weekday BETWEEN 1 AND 7", name="chk_pattern_weekday"),
        CheckConstraint("shift_end > shift_start", name="chk_pattern_time"),
        CheckConstraint("effective_to IS NULL OR effective_to >= effective_from", name="chk_pattern_period"),
        Index("ix_hairdresser_work_patterns_hairdresser_id_weekday", "hairdresser_id", "weekday"),
        Index("ix_hairdresser_work_patterns_salon_id_weekday", "salon_id", "weekday"),
    )

    # Relationships
    hairdresser: Mapped["Hairdressers"] = relationship(back_populates="work_patterns")
    salon: Mapped["Salons"] = relationship(back_populates="work_patterns")
    schedules: Mapped[list["HairdresserSchedules"]] = relationship(back_populates="pattern")