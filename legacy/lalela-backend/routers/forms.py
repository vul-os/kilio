from fastapi_jwt_auth import AuthJWT
from fastapi import Depends
from fastapi import APIRouter

router = APIRouter()

# todo: only get forms for now, no need to create or update them, we can update it manually from the db


# gets list of all forms available
@router.get('/forms')
async def forms(Authorize: AuthJWT = Depends()):
    Authorize.jwt_required()

    current_user = Authorize.get_jwt_subject()
    return {"user": current_user}


# gets the form with the id as a query parameter
@router.get('/forms/{form_id}')
async def form(Authorize: AuthJWT = Depends()):
    Authorize.jwt_required()

    current_user = Authorize.get_jwt_subject()

    return {"user": current_user}
