## My Assumptions

Each player already queries an API every 15 min to see if a new version is available and then updates itself.

Therefore I assumed that the API `PUT /profiles/clientId:{macaddress}`, which the script will call on each machine, has the capacity to get the application's binaries and update itself.

 I also assumed that since players are all different machines, they would be under the same domain and each player's address would be in the form of `https://<server>/profiles/<macaddress>` .

Since we send the same profile to all players and get the same response back on StatusOk, I decided not to test on the response body and test on `StatusCode` instead.

## My Understanding

The tool would read from a .csv file to get the client ids and communicate with each player through this API: `PUT /profiles/clientId:{macaddress}` sending the same body request to all players.

Each player would then have the responsibility to handle the actual update of its application(s). (OUT OF SCOPE)

## My Explanations

I tried to use the standard library as much as possible to limit the amount of dependencies `i.e. net/http, testing`
The only external dependency used is `JWT` for creating a token.

I would normally log errors along the way so we could keep track of the calls that failed during the process and just continue calling the other clients (unless critical error).
It would make it easier to run the update again on those at a later stage

I've created one single auth token with assumption that only one is required for all requests but of course we could have a JWT per player with specific information about the player inside.
This JWT is very basic and is mostly for testing expiration

I would normally stub the inner functions in my unit tests but it's the first time I do unit tests with Golang so it would have required a little more time.

The hard coded `baseURL` (line 59 of main.go) will need to be changed prior to executing the script. It only needs to be done once assuming the baseURL is the same for all clients.

You'll also notice I passed the Args struct as a prefix to the function instead of as arguments. I find it keeps things cleaner by not having everything passed as arguments.