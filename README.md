I did not hand enough time to put this behind SSL

I used own DB storage for simplicity, it stores data in plain files and what is most important provide global locks for writes and reads and transactions which is essential for money transfering.

You can't create new bank, so this list populated manually.

Crypt side should be replaced to use bcrypt algorithm, i just wanted to avoid dependecies


Token based auth, user enter login and password, and receive temporary token. User sessions expiring after 30 minutes, also if new session is started (for example another devide, user will be logged off from first session). Every login gets audited, so you can review it later, however there is no http API for that yet.

Account numbers for the same bank should have same pattern. 

I initially I planned implement multi-bank transfers.