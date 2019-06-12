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

### Common settings

If you need to pass in AWS credentials you can pass them in via the -access-key and -secret-key flags.

The region also needs to be specified. Use -region for this. You only need this if you don't have AWS_REGION specified in your environment. Alternatively EC2_REGION can also be used.

This tool will try to use the default AWS credentials lookups. This includes shared configuration files and EC2 Profiles.

## Help menu

  -access-key string
        Access key for AWS API
  -decrypt
        Request decrypted keys
  -h    Help menu
  -path string
        Parameter Store path
  -recursive
        Look up all keys in branch
  -region string
        Region for AWS API
  -secret-key string
        Secret key for AWS API
  -v    Show Version
