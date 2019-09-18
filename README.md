# go-parameterStoreReader

Simple tool to read data from parameter store

## How to use it

There are 2 ways of collecting keys. Single look ups or path based lookups.

Parameter Store can also store encrypted data which you can request to be decrypted before sending it you.
Use the -decrypt option to enable this.

Use the -path key to specify the path that you would like to view.

All output is sent to STDOUT.

### Single looks ups

Single look ups will only collect a single key and output the data. If the key doesn't exist the tool will error.

### Multiple look ups

Multiple look ups will read a path and then read out all the keys under it.
It will output the keys and values. The key will be displayed then a ":" and then the value.

This is currently under development and will change soon.

This option is used when the -recursive flag is specified.

#### Output

When viewing the output you can have it in 3 formats.

- JSON
- Pretty JSON
- Key:Value

JSON and Pretty JSON is self explaining. It's JSON and 2 space indented JSON.

Key:Value will display the key then a : and then the value. This means you can split on the first ":".
Use in conjunction with --base64 you can split data and not have to worry about new lines making it hard to find the start and end.

#### Nest Keys

If you have nested values in your keys they are displayed in the key name minus the path supplied by the user.

### Common settings

If you need to pass in AWS credentials you can pass them in via the -access-key and -secret-key flags.

The region also needs to be specified. Use -region for this. You only need this if you don't have AWS_REGION specified in your environment. Alternatively EC2_REGION can also be used.

This tool will try to use the default AWS credentials lookups. This includes shared configuration files and EC2 Profiles.

Outputting the data that is collected can be encoded in base64. This is makes it easy to compare values if you want to cross reference values in different places. This is applied only to the values.

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
  -tree-view
        Present output as a tree. Only works with recursive view.
  -v    Show application Version.
```