## Ignite CLI version 

```
$ ignite version
...
Ignite CLI version:	v0.21.0
Ignite CLI build date:	2022-05-10T16:07:28Z
Ignite CLI source hash:	6009368c73266be248d33fe52b9bab653bf17973
Your OS:		darwin
Your arch:		arm64
Your go version:	go version go1.18.2 darwin/arm64
...
```


### Scaffolding NextGame and StoredGame

<mark>The scaffolding of messages does not generate a `creator` field implicity as it occurs in the tutorial.</mark> Example:
```
message NextGame {
    string creator = 1;
    uint64 idValue = 2;
}
```
The tutorial indicates `creator` can be removed from `NextGame` but says nothing about it for `StoredGame`.

<mark>The Genesis type scaffolded is also different from that in the tutorial</mark>.
```
...
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
  NextGame nextGame = 2;
  repeated StoredGame storedGameList = 3 [(gogoproto.nullable) = false];
  // this line is used by starport scaffolding # genesis/proto/state
}
```
### Adjusting default Genesis

This is the DefaultGenesis generated:
```go
// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		NextGame:       nil,
		StoredGameList: []StoredGame{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}
```

And we change to (note the struct of `NextGame` has only one initializer):
```go
// DefaultGenesis returns the default Capability genesis state

func DefaultGenesis() *GenesisState {
	return &GenesisState{
	    NextGame:    &NextGame{uint64(DefaultIndex)},
		StoredGameList: []StoredGame{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}
```

### Including full_game.go

Several changes from the tutorial:
* Now we realize we need `creator` in `StoredGame` so we have to add it to the Protobuf definition.
* Fixing the imports for `rules` that in the repo refers to a different package. This seems to work: `	"github.com/alice/checkers/x/checkers/rules"`  


### Testing

It seems the `genesis_test.go` is already created. 

* In Section `Store Object - make a Checkers Blockchain` the `full_game_test.go ` tests work as expected except for the fact that `TestGetAddressWrongCreator``, TestGetAddressWrongBlack` and `TestGetAddressWrongRed` rely on comparing error message strings that have changed slightly with the new version. For example:
  
`expected: "creator address is invalid: cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4: decoding bech32 failed: checksum failed. Expected 3xn9d3, got 3xn9d4."
                            actual  : "creator address is invalid: cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4: decoding bech32 failed: invalid checksum (expected 3xn9d3 got 3xn9d4)"`
* In the `Message - create a message to create the game` section, the `msg_server_create_game_test.go` uses a function `setupMsgServer` that is scaffolded in the `keeper_test` package, so that test file must be put in `keeper_test` instead of in `keeper` package to work. 




### Commands not working

This gives some error related to addresses (but the tx with no `dry-run` seems to work). 
```
checkersd tx checkers create-game $alice $bob --from $alice --dry-run
```

### Adding the code to message handlers

The code works. 

In `keeper_text.go` there is a call to a function that seems to have changed, concretely:
```go
	keeper := keeper.NewKeeper(
		codec.NewProtoCodec(registry),
		storeKey,
		memStoreKey, paramtypes.Subspace{})  // fix?
```

In the new version the `Keeper` struct has a `paramstore` parameter that is new.
```go
type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
	}
)
```
Also there is a need to prefix `Keeper` with the package name.
```go
func setupKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
...
```

Also, again the tests that are reliant on strings of the error messages have to be fixed:
```
--- FAIL: TestCreateGameRedAddressBad (0.00s)
    msg_server_create_game_test.go:90:
        	Error Trace:	msg_server_create_game_test.go:90
        	Error:      	Not equal:
        	            	expected: "red address is invalid: notanaddress: decoding bech32 failed: invalid index of 1"
        	            	actual  : "red address is invalid: notanaddress: decoding bech32 failed: invalid separator index -1"

        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-red address is invalid: notanaddress: decoding bech32 failed: invalid index of 1
        	            	+red address is invalid: notanaddress: decoding bech32 failed: invalid separator index -1
        	Test:       	TestCreateGameRedAddressBad
FAIL
FAIL	github.com/alice/checkers/x/checkers/keeper	0.403s
FAIL
```

### Adding rejection on games

It seems that after adding the rejection, some of the tests for moving fail since the new `MoveCount` is not correctly handled. These need to be updated at this point, the changes with the correct move counts are in the GitHub repo of the tutorial. 