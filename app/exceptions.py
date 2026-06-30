from fastapi import HTTPException

NotFoundElement=HTTPException(status_code=404, detail="Not found")

NoFieldsToUpdate = HTTPException(status_code=400, detail="No fields to update")