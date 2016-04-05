It defintly took more time then expected, about 7 hours. Mainly because i focused too much on security part, implementing token based auth, but when we talking about financial applications it always pays off.

I used own DB storage for simplicity (i took it from another pet project), it stores data in plain files which simplify debugging a lot and what is most important provide global locks for writes and reads and transactions which is essential for money transfering. In production it should be replaced by DB supporting transactions or by blockchain like Etherium (which i wanted to try first, but i did not had experience with it, so decided to make tradidional way).

App use token based auth, user enter login and password, and receive temporary token, which expire after 30 minutes, or if new session is started (for example another device, user will be logged off from first session). Also there is built in mechanism for deactivating session using flag (for example if account was stolen), but its not exposed to API.

## Completeness: did you complete the features?
Everything from spec is done.
## Correctness: does the functionality act in sensible, thought-out ways?
I try to handle all possible errors, and in human readable way report all the issues. It should suggest approative HTTP codes.
## Maintainability: is it written in a clean, maintainable way?
Definitly codebase is too large for such short amount of time to be nice.., but it's not so bad, especially considering number of tests.
## Testing: is the system adequately tested?
All features covered by tests
## Documentation: is the API well-documented?
I did not added API docs yet, but if i do i would use smth like http://swagger.io/ to do it.


### Things that i did not have time to do:

I did not hand enough time to put this behind SSL

You can't create new bank, so this list populated manually.

Crypt side should be replaced to use bcrypt algorithm, i just wanted to avoid dependecies

Every login gets audited, so you can review it later, however there is no http API for that yet.

Account numbers for the same bank should have same pattern, so we can detect bankID by knowing only account number.

I initially I planned implement multi-bank transfers, but it require more time. But basically i wanted to use smth like etcd for coordination, between banks.

## How to use

#### Creating User:
```
curl -d "name=Leonid Bugaev&dateOfBirth=1986-01-02T15:04:01Z" http://45.33.2.140:8080/user/create
> {"id":"fJt1UIA1oc","password":"y1zhxMnZei"}
```

#### Auth:
```
curl -d "user=fJt1UIA1oc&password=y1zhxMnZei" http://45.33.2.140:8080/auth
{"token":"xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv"}
```

### List accounts
```
curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1oc" http://45.33.2.140:8080/accounts
> {"accounts":[{"Id":"ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb","UserId":"fJt1UIA1oc","Ballance":0}]}
```

### Account Info
```
curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1oc" http://45.33.2.140:8080/account?id=ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb
> {"Id":"ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb","UserId":"fJt1UIA1oc","Ballance":0}
```

### Create Acccount
```
curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1oc" http://45.33.2.140:8080/account/create
> {"id":"bP55j67DL2Hka58TFV4WvH1xrRwzpbXX0c"}
```

### Deposit money
```
curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1oc" -d "to=ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb&amount=1000.1" http://45.33.2.140:8080/deposit
> {}

curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1ochttp://45.33.2.140:8080/account?id=ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb
> {"Id":"ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb","UserId":"fJt1UIA1oc","Ballance":1000.1}
```

### Withdraw money
```
curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1oc" -d "from=ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb&amount=500.1" http://45.33.2.140:8080/withdraw
>{}

curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1ochttp://45.33.2.140:8080/account?id=ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb
> {"Id":"ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb","UserId":"fJt1UIA1oc","Ballance":500}
```

### Transfer money
```
# Create second account (or you can create second user)
curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1oc" http://45.33.2.140:8080/account/create
{"id":"DqyUFkf7JEWphNbJpH9trG31n4BbTQ6ahH"}

curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1oc" -d "from=ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb&to=DqyUFkf7JEWphNbJpH9trG31n4BbTQ6ahH&amount=300" http://45.33.2.140:8080/transfer
>{}

# List accounts
curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1oc" http://45.33.2.140:8080/accounts
> {"accounts":[{"Id":"DqyUFkf7JEWphNbJpH9trG31n4BbTQ6ahH","UserId":"fJt1UIA1oc","Ballance":300},{"Id":"bP55j67DL2Hka58TFV4WvH1xrRwzpbXX0c","UserId":"fJt1UIA1oc","Ballance":0},{"Id":"ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb","UserId":"fJt1UIA1oc","Ballance":200}]

# Amount too large:
curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1oc" -d "from=ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb&to=DqyUFkf7JEWphNbJpH9trG31n4BbTQ6ahH&amount=3000" http://45.33.2.140:8080/transfer
> {"error":"From account does not have enough money"}
```

### View transactions
```
curl -H "X-Auth-Token: xlR2nLipldslXseLM4gE7XEY22FVqg8Uyer86ijmnpg8LPvPSv" -H "X-Auth-User: fJt1UIA1oc" http://45.33.2.140:8080/account/transactions?id=ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb
> [{"Id":"ZbVpvuAFQtzS93rVttz0Gy2Muxzhfh6vEb", ...},...]
```

## Running locally
Repo contains set of docker scripts. If you have docker intalled, first run `make build` to prepare image (only first time). And after you can use `make test` and `make run` commands.

API itself deployed here http://45.33.2.140:8080
