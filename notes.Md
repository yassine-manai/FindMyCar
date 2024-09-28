------------------------------ NOTES ----------------------------
Step 1: 
export PATH=$(go env GOPATH)/bin:$PATH

Step 2:
swag init 

if error exe this :
swag init --parseDependency



StatusMethodNotAllowed - 11 -- 405
StatusNotFound - 9  -- 404
StatusBadRequest  - 12 -- 400
StatusUnauthorized- 401
StatusInternalServerError - 10 - 500
StatusOK - 8 - 200 


StatusMethodNotAllowed - 11
StatusInternalServerError -10
StatusNotFound - 9 
StatusBadRequest  - 12
StatusUnauthorized- 13
StatusOK - 8 

------------------------------- To DO --------------------------
databse/schema creation to take in consideration check it befre everything
autocreate of tables, double check it
Lang Models
List Return for all with id
// map error code
map[int]string


------------------------------- DONE -------------------------
Setup Routes 


































check this i wanna make it as one api that contain extra and if and both of them passed as query and not required and if the extra not passed it return the data without extra elseif extra it return data with extra data and the id query if it passed it return the camerawith id if not passed it retun all data 