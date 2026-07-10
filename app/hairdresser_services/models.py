import uuid

from sqlalchemy import ForeignKey, UUID
from sqlalchemy.orm import mapped_column, relationship, Mapped

from app.database import Base



class HairdresserServices(Base):
    __tablename__ = 'hairdresser_services'

    hairdresser_id: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), ForeignKey("hairdressers.id"),
                                                      primary_key=True)
    service_id: Mapped[int] = mapped_column(ForeignKey("services.id"), primary_key=True)

    # Relationships
    hairdresser: Mapped["Hairdressers"] = relationship(back_populates="hairdresser_services")
    service: Mapped["Services"] = relationship(back_populates="hairdresser_services")