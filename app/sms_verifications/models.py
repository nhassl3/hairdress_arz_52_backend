from datetime import datetime

from sqlalchemy import String, Boolean, DateTime, func, Index, BigInteger
from sqlalchemy.orm import mapped_column, Mapped

from app.database import Base

class SmsVerifications(Base):
    __tablename__ = 'sms_verifications'

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    phone_number: Mapped[str] = mapped_column(String, nullable=False)
    verification_code_hash: Mapped[str] = mapped_column(String, nullable=False)
    expires_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False)
    is_used: Mapped[bool] = mapped_column(Boolean, nullable=False, server_default="false")
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), nullable=False, server_default=func.now())

    __table_args__ = (
        Index("ix_sms_verifications_phone_number_created_at", "phone_number", "created_at"),
    )