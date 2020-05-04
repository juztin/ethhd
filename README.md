# Geth HD Wallet Generator

Creates HD Wallet addresses based on a given password/mnemonic.

### Install

```shell
% go get github.com/juztin/ethhd/cmd/...
```

### Usage

 - Generate and print the key-pairs for 5 accounts of an HD-Wallet  
    ```shell
    % ethhd --accounts 5 --password
    ```
 - Generate and print the key-paris for 5 accounts using a pre-defined mnemonic  
   ```shell
    % ethhd --accounts 5 --mnemonic "...some mnemonic..."
   ```

A password prompt will display for each run, and may be left blank.

Instead of echoing the above commands to the console, you can also store them in a keystore directory, in the
same way Geth stores them, by supplying the `keystore` arg:

```shell
% ethhd \
  --accounts 5 \
  --keystore ./my_keystore
```
