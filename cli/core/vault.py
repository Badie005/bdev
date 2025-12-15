"""
B.DEV CLI - Secrets Vault
Encrypted local secrets management
"""
import json
import base64
import hashlib
from pathlib import Path
from typing import Optional, Dict
from getpass import getpass

# Simple encryption (in production use cryptography library)
def _derive_key(password: str) -> bytes:
    return hashlib.sha256(password.encode()).digest()

def _xor_encrypt(data: str, key: bytes) -> str:
    data_bytes = data.encode('utf-8')
    encrypted = bytes([data_bytes[i] ^ key[i % len(key)] for i in range(len(data_bytes))])
    return base64.b64encode(encrypted).decode('utf-8')

def _xor_decrypt(encrypted: str, key: bytes) -> str:
    encrypted_bytes = base64.b64decode(encrypted.encode('utf-8'))
    decrypted = bytes([encrypted_bytes[i] ^ key[i % len(key)] for i in range(len(encrypted_bytes))])
    return decrypted.decode('utf-8')

VAULT_FILE = Path.home() / "Dev" / ".bdev" / "vault.encrypted"
VAULT_HASH = Path.home() / "Dev" / ".bdev" / ".vault_hash"

class SecretsVault:
    """Encrypted local secrets storage"""
    
    _instance: Optional['SecretsVault'] = None
    _key: Optional[bytes] = None
    _secrets: Dict[str, str] = {}
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
        return cls._instance
    
    def is_initialized(self) -> bool:
        """Check if vault exists"""
        return VAULT_HASH.exists()
    
    def init(self, password: str) -> bool:
        """Initialize new vault"""
        self._key = _derive_key(password)
        self._secrets = {}
        
        # Store password hash for verification
        VAULT_HASH.parent.mkdir(parents=True, exist_ok=True)
        VAULT_HASH.write_text(hashlib.sha256(self._key).hexdigest())
        
        self._save()
        return True
    
    def unlock(self, password: str) -> bool:
        """Unlock existing vault"""
        key = _derive_key(password)
        
        # Verify password
        if not VAULT_HASH.exists():
            return False
        
        stored_hash = VAULT_HASH.read_text().strip()
        if hashlib.sha256(key).hexdigest() != stored_hash:
            return False
        
        self._key = key
        self._load()
        return True
    
    def _load(self):
        """Load secrets from disk"""
        if VAULT_FILE.exists() and self._key:
            try:
                encrypted = VAULT_FILE.read_text()
                decrypted = _xor_decrypt(encrypted, self._key)
                self._secrets = json.loads(decrypted)
            except:
                self._secrets = {}
    
    def _save(self):
        """Save secrets to disk"""
        if self._key:
            data = json.dumps(self._secrets)
            encrypted = _xor_encrypt(data, self._key)
            VAULT_FILE.parent.mkdir(parents=True, exist_ok=True)
            VAULT_FILE.write_text(encrypted)
    
    def set(self, key: str, value: str) -> bool:
        """Set a secret"""
        if not self._key:
            return False
        self._secrets[key] = value
        self._save()
        return True
    
    def get(self, key: str) -> Optional[str]:
        """Get a secret"""
        return self._secrets.get(key)
    
    def delete(self, key: str) -> bool:
        """Delete a secret"""
        if key in self._secrets:
            del self._secrets[key]
            self._save()
            return True
        return False
    
    def list_keys(self) -> list:
        """List all secret keys"""
        return list(self._secrets.keys())
    
    def export_env(self) -> Dict[str, str]:
        """Export secrets as environment variables"""
        return {f"BDEV_{k.upper()}": v for k, v in self._secrets.items()}

def get_vault() -> SecretsVault:
    return SecretsVault()
