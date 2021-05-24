from fastapi_jwt_auth import AuthJWT
from fastapi import Depends
from fastapi import APIRouter

router = APIRouter()


# gets list of all responses
@router.get('/submissions')
async def forms(Authorize: AuthJWT = Depends()):
    Authorize.jwt_required()

    current_user = Authorize.get_jwt_subject()
    return {"user": current_user}


# creates a new submission
@router.post('/submissions')
async def form(Authorize: AuthJWT = Depends()):
    Authorize.jwt_required()

    current_user = Authorize.get_jwt_subject()

    return {"user": current_user}


# updates a created submission
