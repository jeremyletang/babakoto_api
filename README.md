# babakoto_api

## how to build the api

You will need to have go version 1.7 installed on you environment.

Then at the root of the babakoto_api repository just run the following command to retrieve the dependencies:
```shell
$ go get -u -v .
```

Then just run the following command, at the root of the repository again to build the binary:
```shell
$ go build
```

You will need to have the configuration file in the same folder of the binary when you run it.
You can find a configuration example in the repository as well (.babakoto.config.json).

