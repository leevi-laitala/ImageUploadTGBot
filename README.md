# ImageUploadTGBot

Crude telegram bot used to upload images to server the bot is running on.

Created to be simple and fast way to upload images to infoscreen running
on raspberry pi.

In it's current state can only upload images and delete all locally stored
images.

<br>

### Building

Build with your telegram bot API token:

```
$ TOKEN=<token> make build
```

<br>

### Running

```
$ ./citb
```

<br>

### Bot functionality

If sent image as a file, will save the file to specified location on the 
machine on which the bot is running on.

`/delete` command will delete all files in that specified location.

Images sent as just images will be disregarded at the moment.

