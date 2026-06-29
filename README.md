# GATOR blog aggregator


## Installation

To run, you will need Postgres and Go installed

#### Go
Simply run this in your terminal:
>curl -sS https://webi.sh/golang | sh

#### Postgres
Run this in your terminal:
> sudo apt update
> sudo apt install postgresql postgresql-contrib

Update your password:
> sudo passwd postgres

Start the server and enter it:
> sudo service postgresql start
> sudo -u postgres psql

Create the database and connect to it:
> CREATE DATABASE gator;
> \c gator

Set user password:
> ALTER USER postgres PASSWORD 'postgres';

Configure the database by performing migrations. I recommend installing Goose for this:
> go install github.com/pressly/goose/v3/cmd/goose@latest

Then run (you may need to edit this connection string) from the main directory:
> cd sql/schema
> goose postgres postgres://postgres:postgres@localhost:5432/gator up

Download the repo and use "go install ." from inside the main directory to install the gator CLI

Move the .gatorconfig.json file to your home directory. You may need to change the connection string it contains.

## Commands

### Main Commands
**register** - register a new user. Requires name argument
**login** - login to existing user. Requires name argument
**follow** - follows a feed already in the database. Requires url argument
**unfollow** - unfollows a feed. Requires url argument
**addfeed** - adds a new feed and follows it. Requires title and url arguments
**browse** - displays news posts from followed feeds. Optional limit argument to specify how many posts to display (defaults to 2)
**agg** - begin aggregating all feeds in the database. Will continue indefinitely until stopped. Requires a refresh interval argument, formatted as number and unit (10s for 10 seconds, 1h for 1 hour, etc)


### Additional Commands
**reset** - deletes the database. Not reversible!
**users** - displays list of users.
**feeds** - displays list of feeds that can be followed.
**following** - displays list of followed feeds for current user.