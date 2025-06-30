from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from cryptography.hazmat.primitives.ciphers.aead import AESGCM
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.backends import default_backend
import os
import time
import secrets
import logging
import socket

# --- Config ---
TOKEN_EXPIRY_SECONDS = 1800
SESSION_SECRET = secrets.token_bytes(32)
NONCE_SIZE = 12

# --- Logging ---
logging.basicConfig(level=logging.INFO, format='%(asctime)s %(levelname)s %(message)s')
logger = logging.getLogger("nexa-session-manager")

# --- FastAPI app ---
app = FastAPI(title="Nexa Auto Session Manager", version="1.0")

# --- In-memory store ---
token_store = {}

# --- Helper functions ---
def derive_key(secret: bytes) -> bytes:
    """Derives a key from the secret using PBKDF2."""
    kdf = PBKDF2HMAC(
        algorithm=hashes.SHA256(),
        length=32,
        salt=b"nexa-session",
        iterations=100_000,
        backend=default_backend()
    )
    return kdf.derive(secret)

def encrypt_token(token: str, key: bytes) -> dict:
    aesgcm = AESGCM(key)
    nonce = os.urandom(NONCE_SIZE)
    ct = aesgcm.encrypt(nonce, token.encode(), None)
    return {"nonce": nonce.hex(), "ct": ct.hex()}

def decrypt_token(enc: dict, key: bytes) -> str:
    aesgcm = AESGCM(key)
    nonce = bytes.fromhex(enc["nonce"])
    ct = bytes.fromhex(enc["ct"])
    return aesgcm.decrypt(nonce, ct, None).decode()

def sign_token(token: str, timestamp: float) -> str:
    """Generates a HMAC signature for the token and timestamp."""
    import hmac
    sig = hmac.new(SESSION_SECRET, f"{token}:{timestamp}".encode(), digestmod="sha256").hexdigest()
    return sig

# --- Pydantic models ---
class TokenSetRequest(BaseModel):
    token: str

# --- Endpoints ---
@app.post("/set_token")
async def set_token(req: TokenSetRequest):
    token = req.token
    if not token:
        raise HTTPException(status_code=400, detail="No token provided")
    key = derive_key(SESSION_SECRET)
    enc = encrypt_token(token, key)
    timestamp = time.time()
    sig = sign_token(token, timestamp)
    token_store["hf_token"] = {
        "enc": enc,
        "timestamp": timestamp,
        "sig": sig
    }
    logger.info("Token securely stored in session (expires in %ds)", TOKEN_EXPIRY_SECONDS)
    return {"status": "Token stored securely in session.", "expires_in": TOKEN_EXPIRY_SECONDS}

@app.get("/get_token")
async def get_token():
    entry = token_store.get("hf_token")
    if not entry:
        raise HTTPException(status_code=404, detail="No token found")
    now = time.time()
    if now - entry["timestamp"] > TOKEN_EXPIRY_SECONDS:
        token_store.clear()
        logger.warning("Token expired and cleared from session.")
        raise HTTPException(status_code=401, detail="Token expired")
    key = derive_key(SESSION_SECRET)
    token = decrypt_token(entry["enc"], key)
    sig = sign_token(token, entry["timestamp"])
    if sig != entry["sig"]:
        logger.error("Token signature mismatch!")
        raise HTTPException(status_code=403, detail="Token signature mismatch")
    return {"token": token, "expires_in": int(TOKEN_EXPIRY_SECONDS - (now - entry["timestamp"]))}

@app.post("/clear_token")
async def clear_token():
    token_store.clear()
    logger.info("Token cleared from session.")
    return {"status": "Token cleared."}

@app.get("/health")
async def health():
    return {"status": "ok", "uptime": int(time.time())}

@app.exception_handler(Exception)
async def global_exception_handler(exc: Exception):
    logger.error(f"Unhandled error: {exc}")
    return JSONResponse(status_code=500, content={"error": str(exc)})

if __name__ == "__main__":
    logger.info("Starting Nexa Auto Session Manager...")
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        sock.bind(("0.0.0.0", 8765))
        sock.close()
        import uvicorn
        uvicorn.run("session_server:app", host="0.0.0.0", port=8765, reload=False)
    except OSError:
        print("ERROR: Port 8765 is already in use. Please stop the other session_server instance or free the port before starting again.")
