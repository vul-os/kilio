from fastapi_jwt_auth import AuthJWT
from fastapi import HTTPException, Depends
from fastapi import APIRouter
from models.models import User

router = APIRouter()


# provide a method to create access tokens. The create_access_token()
# function is used to actually generate the token to use authorization
# later in endpoint protected
@router.post('/login')
async def login(_user: User, Authorize: AuthJWT = Depends()):
    if _user.username != "test" or _user.password != "test":
        raise HTTPException(status_code=401, detail="Bad username or password")

    # subject identifier for who this token is for example id or username from database
    access_token = Authorize.create_access_token(subject=_user.username)
    return {"access_token": access_token}


# protect endpoint with function jwt_required(), which requires
# a valid access token in the request headers to access.
@router.get('/user')
async def user(Authorize: AuthJWT = Depends()):
    Authorize.jwt_required()

    current_user = Authorize.get_jwt_subject()
    return {"user": current_user}