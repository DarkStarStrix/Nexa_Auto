from flask import Flask, request, jsonify
from cryptography.fernet import Fernet
import os

app = Flask(__name__)

# Generate a key for encryption (in-memory only)
fernet_key = Fernet.generate_key()
cipher_suite = Fernet(fernet_key)

# In-memory encrypted token storage
token_store = {}

@app.route('/set_token', methods=['POST'])
def set_token():
    data = request.json
    token = data.get('token')
    if not token:
        return jsonify({'error': 'No token provided'}), 400
    encrypted_token = cipher_suite.encrypt(token.encode())
    token_store['hf_token'] = encrypted_token
    return jsonify({'status': 'Token stored securely in session.'})

@app.route('/get_token', methods=['GET'])
def get_token():
    encrypted_token = token_store.get('hf_token')
    if not encrypted_token:
        return jsonify({'error': 'No token found'}), 404
    token = cipher_suite.decrypt(encrypted_token).decode()
    return jsonify({'token': token})

@app.route('/clear_token', methods=['POST'])
def clear_token():
    token_store.clear()
    return jsonify({'status': 'Token cleared.'})

if __name__ == '__main__':
    app.run(port=8765)
