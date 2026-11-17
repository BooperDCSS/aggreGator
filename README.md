# aggreGator

aggreGator is an RSS database system for the terminal written in Go and PostgreSQL. It scrapes articles from RSS feeds and stores them in a local database, which you can configure and access using a few simple commands.

It was built following lessons available at [Boot.dev](https://boot.dev).

These instructions assume you are using Linux, macOS, or [WSL](https://learn.microsoft.com/en-us/windows/wsl/install).

# Setup and Installation

This version of aggreGator will work with Go version 1.27.2 or later and PostgreSQL 18.6 or later. First, follow [the instructions for installing Go](https://go.dev/doc/install) on your operating system. Afterward, you will:

- setup the database using PostgreSQL
- create the required JSON config file
- build the binary on your machine
- run commands from the terminal to add, scrape, and review articles

## PostgreSQL Download and Installation

Debian and WSL users should follow the PostgreSQL installation instructions [from Microsoft](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql). 

macOS users can install Postgres with [brew](https://brew.sh/):
```bash
brew install postgresql@18
```
Fedora users will want to follow the installation instructions on [Fedora Docs](https://docs.fedoraproject.org/en-US/quick-docs/postgresql/).

Check that you downloaded and successfully installed the correct version by running `psql --version` from your terminal.

On Linux/WSL, you need to set a system-level password:
```bash
sudo passwd postgres
```
Make sure you don't forget it.

## PostgreSQL Setup
Once installed, start the PostgreSQL server:
- on Linux/WSL: `sudo service postgresql start`
- on Mac: `brew services start postgresql@18`

Connect to the server using `psql`, which comes bundled with your PostgreSQL installation:
- on Linux/WSL: `sudo -u postgres psql`
- on Mac: `psql postgres`

Create a new database called 'gator'. These instructions assume that is the name you will use. If you decide on a different name, replace all uses of the term 'gator' with your choice: 
```sql
CREATE DATABASE gator;
```
Then connect to the database with `\c gator`. If you see a `gator=#` prompt, the server is ready.

On Linux/WSL, you will need to set a user password in addition to your system-level password. Since the server runs locally, keep it simple. From the `gator=#` prompt:
```sql
ALTER USER postgres PASSWORD 'postgres';
```

Type `exit` to leave the `psql` shell. Create a JSON config file at `~/.gatorconfig.json` and paste the following into it according to your OS:

- Linux/WSL:
```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
}
```

- macOS (replace 'your_user_name' with the username you use to login):
```json
{
  "db_url": "postgres://your_user_name:@localhost:5432/gator?sslmode=disable"
}
```

Test your connection, e.g. on Linux:
```bash
psql 'postgres://postgres:postgres@localhost:5432/gator'
```
If you connect to the server and see the `gator=#` prompt, everything is ready.

## Clone the Repository

Choose the directoy where you would like to store the program. I store mine in `~/src`, but `~/repos` or any other folder works. From there, clone the repo:
```bash
git clone https://github.com/BooperDCSS/aggreGator
```
This will create a folder for the project, e.g., `~/src/aggreGator`.

## Migrate the Database

Install [Goose](https://github.com/pressly/goose):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
Run `goose -version` from your terminal to ensure installation was successful.

Change your directory to `aggreGator/sql/schema`. From here, migrate the database using Goose:
```bash
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
```
You should see the message `goose: successfully migrated database to version: 5`. We are now ready to use aggreGator.

# Installing and Using Aggregator

You should `go build` or `go install` aggreGator from the project's root directory. The former will build an executable file in the project's root directory, the latter will build an executable in the GOBIN directory, which means you can run it from anywhere in your system. By default, the name of the resulting program will be `aggregator`. You can change the name of your directory to something else if you want the program to have a different name, like `gator`. I recommend installing it, but you do you. If you don't feel like compiling, you can use `go run .` in place of `aggregator` in the following instructions.

# Commands

All commands follow the pattern `aggregator <command name> <sometimes optional arguments>`. If arguments are required or available for a command, they are listed under the command name below.

## register

aggreGator uses a passwordless user system, which allows different users to store different RSS feeds. In practice, you can use this to manage lists of feeds by topic or interest. At least one user must exist in the database before you can begin managing feeds. The `register` command both creates and sets a username to the current user in the database:

```bash
aggregator register yourNameHere
# or use your full name
aggregator register "Your Full Name with Spaces"
```

## login

If you have registered multiple names, you can switch betweeen them with `login`:
```bash
aggregator login yourNameHere
```

## reset

`aggregator reset` will erase all of the users and all of their feeds from the database. This also destroys all scraped articles. Use this if you want a fresh start.

## users

`aggregator users` lists all all registered users in the database and indicates which of them is the current user.

## addfeed

Requires "feed name" and "URL" arguments. This is how you add a feed to your list:
```bash
aggregator addfeed "This Week in Linux" "https://feeds.fireside.fm/thisweekinlinux/rss"
```

## agg

Requires a "frequency" argument. Scrapes your saved feeds on a loop that runs according to the defined frequency. Begins with unscraped feeds and then proceeds according to "last updated" data, beginning with the least up-to-date feeds:
```bash
aggregator agg 30s # scrapes every 30 seconds
aggregator agg 1m # every minute
```

The program will accept many time units, including 's' for seconds, 'm' for minutes, and 'h' for hours, but I strongly recommend using `30s` or more to avoid flooding your favorite feeds with requests. You can stop the loop with `CTRL+C`. The program will not scrape posts already stored in the database. You will be informed when that starts to happen:

```bash
"NOTICE ---> URL already in database, skipping post."
```

## browse

Accepts an optional "count" argument. Allows you to browse saved articles from your feeds starting with the most recently published ones (this uses published data from the feeds themselves). Without the optional argument, it will return the two most recent articles.

```bash
aggregator browse 10 # returns the 10 most recently scraped articles
```
The "count" argument describes an upper bound. If you request 300 articles but only have 30 saved, aggreGator will only return 30 articles.

## feeds

`aggregator feeds` lists all stored feeds, their URLs, and the user that created them.

## follow

Requires a complete URL argument. If a feed already exists in the database, aggreGator won't let you add it a second time. If you want to follow the same feed across two users, run `aggregator follow <url>`. This will add an already existing feed to the current user's list of feeds.

## following

`aggregator following` lists the names of all the feeds followed by the current user.

## unfollow

Requires a complete URL argument. `aggregator unfollow <url>` does what it says on the tin: unfollows a feed. It will also remove all of the associated posts for that feed from your database.

