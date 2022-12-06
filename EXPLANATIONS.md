## My Assumptions

Each player already queries an API every 15 min to see if a new version is available and then updates itself.

Therefore I assume that the API `PUT /profiles/clientId:{macaddress}`, which the tool will call on each machine, has the capacity to get the application's binaries and update itself.

 I assume that since players are all different machines, they would be under the same domain and each player's address would be in the form of ` https://<server>/profiles/<macaddress>` .

Since we send the same profile to all players and get the same response back on StatusOk, there's no point testing on the response body.
That's why I decided to test on `StatusCode` instead.

## My Understanding

The tool would read from a .csv file to get the client ids and communicate with each player through this API: `PUT /profiles/clientId:{macaddress}` sending the same body request to all players.

Each player would then have the responsibility to handle the actual update of its application(s). (OUT OF SCOPE)

## My Explanations

I tried to use the standard library as much as possible to limit the amount of dependencies `i.e. net/http, testing`
The only external dependencies used are `JWT` for creating a token and `gostub` so it's easier/faster to implement tests for me at this point but I would rather do without.

I am passing the arguments `baseURL, clientId, token` to functions 'cause that way it's easy to replace when testing

I would normally log any errors along the way but continue calling other players (unless critical error) since there might be a lot of players to update. This way I would know which one failed and be able to run the update again on those.

I've created one single auth token with assumption that only one is required for all updates but of course we could have a JWT per player with specific information about the player inside.
This JWT is very basic and is mostly for testing expiration

I would normally stub the inner functions in my unit tests but it's the first time I do unit tests with Golang.