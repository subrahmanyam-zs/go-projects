This is the example to understand the read and write functionality for different filestores(AWS,AZURE,GCP,SFTP,LOCAL).
Just set the `FILE_STORE` value to your desired file store type and set other configs accordingly.

For reading a file , you should have a `test.txt` file in your filestore which will be containing the content to be read.
By default it will read only 20 bytes, you need to set the number of bytes you want to read, in handler.