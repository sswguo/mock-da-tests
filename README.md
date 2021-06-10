# query-metadata

The script is used to grab the alignment logs from PNC
and extract the gavs from it, then request the metadata 
of the gavs from Indy concurrently.

# config
- pnc_rest_url: PNC Rest url
- indy_url: Indy url
- da_group: The name of the maven group repository used for DA.
- max_concurrent_goroutines: The number of the concurrent threads to access Indy.

# parameters
- buildId: the build id generated in PNC.
