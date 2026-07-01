from fastapi import HTTPException

NotFoundElement=HTTPException(status_code=404, detail="Not found")

NoFieldsToUpdate = HTTPException(status_code=400, detail="No fields to update")

AlreadyExistsElement=HTTPException(status_code=409, detail="Already exists")

UserHasBookings = HTTPException(status_code=409, detail="User has bookings")