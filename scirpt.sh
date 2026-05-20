curl -X POST http://localhost:8080/v1/auth/register \
-d '{
    "password":"Arsen",
    "email":"transformatoratagac@gmail.com"
}'

{"access":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzkwNTUxMDgsInJvbGUiOiJzdXBlci1hZG1pbiIsInN1YiI6N30.NkoL3U2Nu-UybUNPvQ_-RStzuIQf962wzwyOyo0GX1k","refresh":"914e8a4bf6939d76ef3f7e3d96ef88624cd6bd04fd879c8f5adb5bac3e796253"}

curl -X POST http://localhost:8080/v1/auth/refresh \
    -d '{"token":"914e8a4bf6939d76ef3f7e3d96ef88624cd6bd04fd879c8f5adb5bac3e796253"}'

curl -X PATCH http://localhost:8080/v1/products/12/ \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzkwNTQwNTMsInJvbGUiOiJzdXBlci1hZG1pbiIsInN1YiI6N30.Vmtn4agIDacDgtWmkUWp3goo5WjnYLXm6SjzMYeV_8w"\
    -d '{
        "name":"bir"
    }'

curl -X PATCH http://localhost:8080/v1/products/12/image \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzkwNTQwNTMsInJvbGUiOiJzdXBlci1hZG1pbiIsInN1YiI6N30.Vmtn4agIDacDgtWmkUWp3goo5WjnYLXm6SjzMYeV_8w"\
    -F "image=@/home/chupep/images/Lain.jpg"

curl -X HEAD http://localhost:8080/v1/password \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzkwMjk0NTksInJvbGUiOiJ1c2VyIiwic3ViIjo3fQ.Gc1mwW9ZUZYomjuwmXKm6kB0VX0oH3-bLln-E4-Ljn4"\
    -d '{"email":"transformatoratagac@gmail.com"}'

curl -X POST http://localhost:8080/v1/password/42453 \
  -d '{
    "password":"Arseniy",
    "email":"transformatoratagac@gmail.com"
}'

curl -X GET http://localhost:8080/v1/admin \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzkwNTQwNTMsInJvbGUiOiJzdXBlci1hZG1pbiIsInN1YiI6N30.Vmtn4agIDacDgtWmkUWp3goo5WjnYLXm6SjzMYeV_8w"\

curl -X GET http://localhost:8080/v1/admin?page=1&limit=10 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzkwNTUxMDgsInJvbGUiOiJzdXBlci1hZG1pbiIsInN1YiI6N30.NkoL3U2Nu-UybUNPvQ_-RStzuIQf962wzwyOyo0GX1k"\

curl -X PATCH http://localhost:8080/v1/admin/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzkwNTUxMDgsInJvbGUiOiJzdXBlci1hZG1pbiIsInN1YiI6N30.NkoL3U2Nu-UybUNPvQ_-RStzuIQf962wzwyOyo0GX1k"\
  -d '{"role":"loh"}'

{"access":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzkwNjA0OTYsInJvbGUiOiJzdXBlci1hZG1pbiIsInN1YiI6OH0.xLRLmLNt3DsbZtwQPFN3zWr-SS6ZIqprRtPj_7TUkik","refresh":"bfac2df364e9e5100634180a20338b8c4ad41f20c263294e5b7eaa789e8612ff"}

{"access":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzkyMjk0NTgsInJvbGUiOiJzdXBlci1hZG1pbiIsInN1YiI6OH0.dHidY2JdFTFiMec7JVBI4HVFAYhaHDy5e0uKigdPcJ0","refresh":"f2572a3dbb2413a95d1c4dc0bccc14324a3b16d5687a4c9c6a0192a18252b594"}