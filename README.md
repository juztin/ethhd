# Geth HD Wallet Generator

Creates HD Wallet addresses based on a given password/mnemonic.

### Install

```shell
% go install github.com/juztin/gethhdwallet/cmd/gethhdgenerator
```

### Usage

 - Generate and print the key-pairs for 5 accounts of an HD-Wallet  
    ```shell
    % gethhdgenerator --accounts 5 --password "super_secret"
    ```
 - Generate and print the key-paris for 5 accounts using a pre-defined mnemonic  
   ```shell
    % gethhdgenerator --accounts 5 --mnemonic "...some mnemonic..."
   ```
 - Generate and print the key-paris for 5 accounts using a pre-defined mnemonic and password  
   ```shell
    % gethhdgenerator --accounts 5 --password "super_secret" --mnemonic "...some mnemonic..."
   ```

Instead of echoing the above commands to the console, you can also store them in a keystore, in the
same format the Geth stores them, by supplying the `keystoredir` and `keystorepassword` args:

```shell
% gethhdgenerator \
  --accounts 5 \
  --password "super_secret" \
  --keystoredir ./my_keystore \
  --keystorepassword "super_secret"
```
