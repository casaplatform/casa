Random notes for now...

Casa is designed to be as flexible as possible. That is why this package is
nothing but type definitions, and everything that Casa can do is done through 
external service "plugins". 

The cmd/casa sub-package is the 
default server, and it uses three plugins right now: mqtt, logger, and hue. When
run it starts an MQTT broker, uses the logger to print out all received
messages, and starts the hue service to connect to a hue bridge. You'll need to
know the bridge IP and have a user name already generated. This will be
automated in the future. 


* Casa starts an MQTT broker on port 1883 with ZERO security. This should
  probably be the first improvement.

* We should also add the ability to use an external MQTT broker if desired.

* Use the `github.com/pkg/errors` package for returning errors where possible. 
  It's very handy and provides a nice stack trace if we want it.

* `github.com/spf13/viper` is used for handling configuration. Each service has 
  it's own section in the config file, and that section is passed to the service
  when it is started. 

* `github.com/spf13/cobra` is used to handle the casa program commands. Right now
  there is just a command to start the server and print the version. More should
  be added, like the ability to generate a default config file on demand.

* There is an `environment` package that `cmd/casa` uses to store global
  configuration and state. This makes using a single logger easier, as well as 
  giving somewhere the plugins can register themselves to. 
