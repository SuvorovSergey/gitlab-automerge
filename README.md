# gitlab-automerge

The bot will track all merge reqests that are not marked as Draft and automatically accept them if the number of likes reaches 2 (or other value that is specified in the settings)

## installation

clone repository
```
git clone git@github.com:SuvorovSergey/gitlab-automerge.git
```

then you need to install the dependencies and create a configuration file

```
go mod vendor
cp config.yml.example config.yml

```

or use Makefile command:

```
make init
```



Run app:
```
go run ./cmd/main.go
```

or use Makefile command:

```
make run
```

Build:
```
GOOS=linux go build -o ./bin/gitlab-automerge ./cmd/main.go
```
or use Makefile command:

```
make build
```

---
## Gitlab project configuration

In each project in gitlab, which should be tracked by the bot, you must add a configuration file to the default branch

Example: 
filename ".automerge"
```
{
  "upvotes_threshold": 2 
}

```

Then you need to create the user under which the bot will work and add the user to the project as a maintainer. Write the user's private key in the config.yml configuration file

The bot will track all merge reqests that are not marked as Draft and automatically accept them if the number of likes reaches 2 (or other value that is specified in the settings)



