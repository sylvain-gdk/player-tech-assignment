## My Assumptions

Each player already queries an API every 15 min to see if a new version is available and then updates itself.

Therefore I will assume that the API `PUT /profiles/clientId:{macaddress}`, which the tool will call on each machine, has the capacity to get the application's binaries and update itself.

Since we send the same profile to all players and get the same response back, there's no point testing on the response body.

## My Understanding

The tool would read from a .csv file to get the client ids and communicate with each player through this API: `PUT /profiles/clientId:{macaddress}` sending the same body request to all players.

Each player would then have the responsibility to handle the actual update of its application(s). (OUT OF SCOPE)

## My Explanations

I tried to use the standard library as much as possible to limit the amount of dependencies i.e. net/http, test