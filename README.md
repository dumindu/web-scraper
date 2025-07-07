# web-scraper

## Database Design

![Database Design](database/Design.png)

## Email Design

![Email Design](internal/mailer/assets/Design.png)

Credits go to [leemunroe/responsive-html-email-template](https://github.com/leemunroe/responsive-html-email-template)

> MailHog: http://localhost:8025
> Username: `user` Password: `password`

## Ed25519 JWT Signing and Verification Keys

![Token Design](internal/utils/jwtutil/assets/Token.png)

> [Public/ Private PEM keys for access/refresh tokens](internal/utils/jwtutil/assets)
> 
> ```bash
> openssl genpkey -algorithm ed25519 -outform PEM -out access-private-key.pem
> openssl pkey -outform PEM -pubout -in access-private-key.pem -out access-public-key.pem
> ```
