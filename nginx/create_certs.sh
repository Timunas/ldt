#! /bin/bash
echo """
==================================================
  Creating certificate for timunas.test.dev
==================================================
"""
mkcert timunas.test.dev

echo """
==================================================
  Copying timunas.test certificates to ./nginx/certs
==================================================
"""

mv ./timunas.test.dev.pem ./nginx/certs/timunas.test.dev.crt
mv ./timunas.test.dev-key.pem ./nginx/certs/timunas.test.dev.key

echo """
==================================================
  Add the following to /etc/hosts file:
  127.0.0.1 timunas.test.dev
==================================================
"""
