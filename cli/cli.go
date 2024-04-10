package cli

type CLI struct {
	Rows     int    `help:"Number of processed rows" default:"1000000"  short:"r"`
	Datafile string `help:"Name of the datafile containing JSON data"  short:"d" required:""`
	Postgres struct {
		Host     string `help:"Set the hostname"  short:"H" default:"localhost"`
		Port     int    `help:"Set listening port"  short:"P" default:"5432"`
		User     string `help:"Set the database username" short:"U"`
		Password string `help:"Set the database password"  short:"W"`
		Database string `help:"Set the database name"  short:"D"`
	} `cmd:"" help:"Execute PostgreSQL tests"`
	Elasticsearch struct {
		Host  string `help:"Set the hostname"  short:"H" default:"localhost"`
		Port  int    `help:"Set listening port"  short:"P" default:"9200"`
		Index string `help:"Set the index name"  short:"I"`
	} `cmd:"" help:"Execute Elasticsearch tests"`
	MySQL struct {
		Host     string `help:"Set the hostname"  short:"H" default:"localhost"`
		Port     int    `help:"Set listening port"  short:"P" default:"3306"`
		User     string `help:"Set the database username" short:"U"`
		Password string `help:"Set the database password"  short:"W"`
		Database string `help:"Set the database name"  short:"D"`
	} `cmd:"" help:"Execute MySQL tests" name:"mysql"`
	MongoDB struct {
		Host string `help:"Set the hostname"  short:"H" default:"localhost"`
		Port int    `help:"Set listening port"  short:"P" default:"27017"`
		// User     string `help:"Set the username" short:"U"`
		// Password string `help:"Set the password" short:"W"`
		Database string `help:"Set the database name"  short:"D"`
	} `cmd:"" help:"Execute MongoDB tests" name:"mongo"`
	CouchDB struct {
		Host     string `help:"Set the hostname"  short:"H" default:"localhost"`
		Port     int    `help:"Set listening port"  short:"P" default:"5984"`
		User     string `help:"Set the username" short:"U"`
		Password string `help:"Set the password" short:"W"`
		Database string `help:"Set the database name"  short:"D"`
	} `cmd:"" help:"Execute CouchDB tests" name:"couchdb"`
}
