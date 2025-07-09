from fastapi import FastAPI, Request, HTTPException
from fastapi.responses import JSONResponse
import uvicorn

app = FastAPI()

SECRET_HEADER = "X-App-Token"
SECRET_VALUE = "your_very_secret_token_12345"


@app.middleware("http")
async def check_secret_header(request: Request, call_next):
    if request.headers.get(SECRET_HEADER) != SECRET_VALUE:
        return JSONResponse(
            status_code=403, content={"detail": "Forbidden: Invalid token"}
        )
    return await call_next(request)


@app.get("/greet")
async def greet(name: str):
    return {"message": f"Привет {name}"}


@app.get("/health")
async def health_check():
    return {"status": "ok"}


if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=8000, log_level="info")
