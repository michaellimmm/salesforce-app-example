# Salesforce app example

This is example of salesforce connected app

```mermaid
sequenceDiagram
    autonumber

    participant user as User
    participant app as App
    participant sf as Salesforce

    user->>+app: [REST] /auth
    app->>+sf: Redirect url
    sf->>-app: Return code
    app->>+sf: get authorization_code
    sf->>-app: return
    app->>-user: Response
```

```
https://connect-momentum-3503.my.salesforce.com/services/oauth2/authorize?
client_id=3MVG9fe4g9fhX0E6ZWCD0XNigpXEN5swRzlbInje3qHqoA4z0rVY0gsHEPQZl1bYTIblOmGF40OLWm5ylqvTl&
redirect_uri=https://webhook.site/0e9734d9-def1-4ac9-86cc-311f3be074ac&
response_type=code
```
