from datetime import datetime

from pydantic import BaseModel


class SendSmsVerification(BaseModel):

    phone_number: str


class VerifySmsVerification(BaseModel):
    phone_number: str
    verification_code_hash:str


class SmsVerificationResponse(BaseModel):
    id:int
    phone_number: str
    verification_code_hash:str
    expires_at: datetime
    is_used: bool
    created_at: datetime