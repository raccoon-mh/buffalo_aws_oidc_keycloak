# buffalo_aws_oidc_keycloak

## info

- golang Buffalo 프레임워크를 기반으로 Keycloak OIDC 클라이언트를 사용해 AWS IAM 자격증명 공급자(IDP)로써 STS 토큰을 발급받아 VM List 를 호출할 수 있는 테스트.
- http 프로토콜 사용으로 단순 ReadOnly 테스트 및 학습 용도로써 사용 권고.

## Layout
![Untitled (11)](https://github.com/raccoon-mh/buffalo_aws_oidc_keycloak/assets/130422754/5d73d5fb-3e59-4f9d-8528-ad3d61ceabb1)

# env

- 본 테스트를 사용하기 위해 아래 기본 /.env 가 필요함.

```
ADDR = 0.0.0.0 
PORT = 4000

keycloakHost = {keycloakHost} # mykeycloak.domain.com
realm = {realm} # Realm Name
client = {client} # Client Name (Aud)
clientSecret = {clientSecret} # client Secret
RoleArn = arn:aws:iam::{ACCT ID}:role/{Role NAME}
```

## How To USE (Tested Ubuntu 20.04)

1. install Keycloak and set Realm, Aws Client…
    1. wiki…
2. Set Aws IAM IDP & Role…
    1. wiki…
3. install golang and buffalo(+node)
    1. https://go.dev/doc/install [install go (go version go1.21.4 linux/amd64)]
    2. https://gobuffalo.io/documentation/getting_started/installation/ [install go buffalo (INFO[0000] Buffalo version is: v0.18.8)]
    3. install node (tested v20.5.1). recommand use nvm (https://github.com/nvm-sh/nvm)
4. clone this repo

```bash
git clone https://github.com/raccoon-mh/buffalo_aws_oidc_keycloak 
```

5. make .env file for  buffalo

```bash
cd buffalo_aws_oidc_keycloak/
nano ./.env
```

6. start buffalo dev


```bash
cd buffalo_aws_oidc_keycloak/
buffalo dev
```

## Home(http://localhost:4000/)
![Untitled (6)](https://github.com/raccoon-mh/buffalo_aws_oidc_keycloak/assets/130422754/ab166444-3900-4704-83cc-c7d977de43e0)


## Login(http://localhost:4000/login/)
![Untitled (7)](https://github.com/raccoon-mh/buffalo_aws_oidc_keycloak/assets/130422754/bd656689-2110-48f6-861f-3f08f9ca339a)
- Keycloak 계정으로 로그인

## User Home(http://localhost:4000/user/home)
![Untitled (8)](https://github.com/raccoon-mh/buffalo_aws_oidc_keycloak/assets/130422754/85c8e127-4759-446f-bb8e-0b29742f5be6)
![Untitled (9)](https://github.com/raccoon-mh/buffalo_aws_oidc_keycloak/assets/130422754/2f10e5dd-5226-4eeb-8b01-e8aff21467ae)
- GET STS 버튼을 통해 AWS STS 토큰 발급 가능

![Untitled (10)](https://github.com/raccoon-mh/buffalo_aws_oidc_keycloak/assets/130422754/b1b98225-a358-4f8a-93b5-5088b91d5889)
- VM List 호출 가능
