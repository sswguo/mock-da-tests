# Mock-Dependency-Analysis

The script is used to grab the alignment logs from PNC
and extract the gavs from it, then request the metadata 
of the gavs from Indy concurrently.

# parameters
- `pnc_rest_url`: PNC Rest url
- `indy_url`: Indy url
- `da_group`: The name of the maven group repository used for DA.
- `max_concurrent_goroutines`: The number of the concurrent threads to access Indy.
- `buildId`: the build id generated in PNC.

# run as container
```
docker build -t quay.io/wguo/mockda .

docker run  --env PNC_REST=http://<ORCH_HOST>/pnc-rest/v2 --env INDY_URL=http://<INDY_HOST> --env DA_GROUP=DA --env BUILD_ID=<BUILD_ID> quay.io/wguo/mockda
```
