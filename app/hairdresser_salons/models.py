import uuid

from sqlalchemy import  ForeignKey,  UUID
from sqlalchemy.orm import mapped_column, Mapped, relationship


from app.database import Base

class HairdresserSalons(Base):
    __tablename__ = 'hairdresser_salons'

    hairdresser_id: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), ForeignKey("hairdressers.id"),
                                                      primary_key=True)
    salon_id: Mapped[int] = mapped_column(ForeignKey("salons.id"), primary_key=True)

    # Relationships
    hairdresser: Mapped["Hairdressers"] = relationship(back_populates="hairdresser_salons")
    salon: Mapped["Salons"] = relationship(back_populates="hairdresser_salons")