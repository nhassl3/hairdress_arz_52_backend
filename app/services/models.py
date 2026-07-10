from datetime import timedelta
from decimal import Decimal

from sqlalchemy import String,  Text, Interval, Numeric, CheckConstraint
from sqlalchemy.orm import mapped_column, relationship, Mapped

from app.database import Base


class Services(Base):
    __tablename__ = 'services'

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    service_name: Mapped[str] = mapped_column(String, nullable=False)
    duration: Mapped[timedelta] = mapped_column(Interval, nullable=False)
    price: Mapped[Decimal] = mapped_column(Numeric(10, 2), nullable=False)
    description: Mapped[str] = mapped_column(Text, server_default="")

    __table_args__ = (
        CheckConstraint("price > 0", name="check_price_positive"),
    )

    # Relationships
    hairdresser_services: Mapped[list["HairdresserServices"]] = relationship(back_populates="service",
                                                                            cascade="all, delete-orphan")
    bookings: Mapped[list["Bookings"]] = relationship(back_populates="service")