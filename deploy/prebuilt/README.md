# Prebuilt EVM configurations

This is currently a very basic implementation to copy configuration files into place. Only files that need to be overwritten are included. 

To use cd to the **/deploy** directory and run:
```bash
make prebuilt PREBUILT=babblepoa
``` 

There is a new folder prebuilt under the deploy directory. Within that folder is a templates folder that contains the files to overwrite the current configuration. The production release will use zipped files rather than the file system.

The name of the directory under the templates folder is the name of the PREBUILT instance. 

In the root of the instance may be an optional .message file that is displayed to the user when that install that template. This would typically contain the make start command parameters for this instance. 

Current instances are

- **babblepoa** - Babble 4 node instance with POA smart contract in the Genesis block 