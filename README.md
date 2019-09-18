# go-parameterStoreReader

Simple tool to read data from parameter store

## How to use it

There are 2 ways of collecting keys from AWS Parameter Store. Single look ups or path based lookups.

Parameter Store can store encrypted data which you can request to be decrypted before sending it you.
Use the -decrypt option to enable this.

Use the -path key to specify the path that you would like to view.

All output is sent to STDOUT by default or its written to a file if the `-f` flag is given.

### Single looks ups

Single look ups will only collect a single key and output the data.

If the key doesn't exist or is not the end of a branch the tool will error.

### Recursive look ups

Recursive look ups will walk a path and then read out all the keys under it.
It will output the keys and values. The key will be displayed using the format requested.

This option is used when the -recursive flag is specified.

#### Output

Output can be given in the following formats.

- json
- pretty-json
- line
- env

By default the path supplied is removed from the output. It can be left in place by using the `-include-path` flag.

##### json

The output is given in JSON

##### pretty-json

The output is given with 2 space indented JSON. Human readable JSON is what we are going for here.

##### yaml

The output is given in YAML format

##### line

The output will display the key then a `:` and then the value. This means you can split on the first ":".
Each value will be presented on a new line.

```text
key1/key2/final1:value1
key1/key2/final2:value2
key1/keyABC/final1:valueABC
```

Pro Tips:

* Remember that the presented path is stripped of the passed in path unless you use the `-include-path` flag.
* Use in conjunction with --base64 you can split data and not have to worry about new lines making it hard to find the start and end.

##### env

The output is best used when setting environment values from Parameter Store.

```text
KEY123=VALUE-ABC
KEY456=VALUE-XYZ
```

The output is stripped of the path regardless of `-include-path`.

### Common settings

If you need to pass in AWS credentials you can pass them in via the `-access-key` and `-secret-key` flags.

The region for which Parameter Store you want to connect to needs to be specified. Use the `-region` for this.
If no region is specified `AWS_REGION` or `EC2_REGION` specified in your environment will be used.

This tool will try to use the default AWS credentials lookups. This includes shared configuration files and EC2 Profiles. Use `-profile` if you want to specify a profile.

Data that is collected can be encoded in base64 using the `-base64` flag. This removed all new line markers making it easy to compare values if you want to cross reference values in different places. This is applied only to the values.

## Examples

```sh
# Single collection with base64 output
./ps-reader -base64 -path /app1/DB_USERNAME

# Read all values under /app1 and return them as pretty json
./ps-reader -base64 -format pretty-json -include-path -recursive -path /app1

# Set environment variables from Parameter store
for ev in $(./ps-reader -format env -region eu-west-1 -decrypt -recursive -path /c3/infrastructure/telegraf); do
    eval "export $ev"
done
```

## Help menu

```
  -access-key string
        Access key for AWS API.
  -base64
        Base64 encode collected values.
  -config-file string
        AWS Config file override, only valid with -profile.
  -decrypt
        Request decrypted keys.
  -f string
        Output to specified file.
  -format string
        Format for output. Supported values: line,json,pretty-json,yaml. (default "line")
  -h    Help menu.
  -include-path
        Include the passed in path in the output. Only used with recursive lookups.
  -path string
        Parameter Store path.
  -profile string
        AWS Profile to use.
  -recursive
        Look up all keys in branch.
  -region string
        Region for AWS API.
  -secret-key string
        Secret key for AWS API.
  -v    Show application Version.
```
