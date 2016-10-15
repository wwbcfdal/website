cat wwbc-aws-key-pair.pem.enc | openssl enc -d -aes-256-cbc > wwbc-aws-key-pair.pem.dec
cat            classified.enc | openssl enc -d -aes-256-cbc > classified.dec
rm wwbc-aws-key-pair.pem.enc
rm classified.enc
